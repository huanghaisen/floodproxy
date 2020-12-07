package loadbalance

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/g4zhuj/hashring"
)

// 多种选择方式
const (
	SelectModeRandom int = 1 //随机选择
	SelectModePoll   int = 2 //轮询选择
	SelectModeHash   int = 3 //哈希选择
)

// HostInfo Host 地址
type HostInfo struct {
	Target          string             //转发目标域名
	MultiTarget     []string           //目标域名
	IsMultiTarget   bool               //是否有多转发目标
	MultiTargetMode int                //多转发目标选择模式
	PoolModeIndex   int                //轮询模式索引
	hashRing        *hashring.HashRing //一致性哈希
}

// HostInfoInterface interface for host
type HostInfoInterface interface {
	GetTarget(req *http.Request) string
}

var ipForwardeds []string

// 如果消息是通过前端代理服务器转发或者cdn转发，则需要从消息头中获取IP地址（注意确保IP的真实性），
// 如果消息直接来自于用户客户端，则使用req.RemoteAddr获取
// getIPAddr 获取IP 地址
func getIPAddr(req *http.Request) []string {
	if ipForwardeds == nil {
		return []string{strings.Split(req.RemoteAddr, ":")[0]}
	} else {
		for _, v := range ipForwardeds {
			if addr, ok := req.Header[v]; ok && len(addr) > 0 {
				return addr
			}
		}
		return []string{strings.Split(req.RemoteAddr, ":")[0]}
	}
}

// GetTarget 选取 目标
func (hostInfo *HostInfo) GetTarget(req *http.Request) string {
	var route string
	if hostInfo.IsMultiTarget {
		if hostInfo.MultiTargetMode == SelectModeRandom { //随机模式
			route = hostInfo.MultiTarget[rand.Int()%len(hostInfo.MultiTarget)]
		} else if hostInfo.MultiTargetMode == SelectModePoll { //轮询模式
			route = hostInfo.MultiTarget[hostInfo.PoolModeIndex]
			hostInfo.PoolModeIndex++
			hostInfo.PoolModeIndex = hostInfo.PoolModeIndex % len(hostInfo.MultiTarget)
		} else if hostInfo.MultiTargetMode == SelectModeHash { //哈希模式
			ips := getIPAddr(req)
			route = hostInfo.hashRing.GetNode(ips[0])
		} else { //未配置或配置错误使用随机模式
			route = hostInfo.MultiTarget[rand.Int()%len(hostInfo.MultiTarget)]
		}
	} else {
		route = hostInfo.Target
	}
	return route
}

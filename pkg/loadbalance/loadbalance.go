package loadbalance

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/serialx/hashring"
)

// 多种选择方式
const (
	SelectModeRandom int = 1 //随机选择
	SelectModePoll   int = 2 //轮询选择
	SelectModeHash   int = 3 //轮询选择
)

// HostInfo Host 地址
type HostInfo struct {
	Target          string             //转发目标域名
	MultiTarget     []string           //目标域名
	Length          int                //可转发的数量
	IsMultiTarget   bool               //是否有多转发目标
	MultiTargetMode int                //多转发目标选择模式
	PoolModeIndex   int                //轮询模式索引
	hashRing        *hashring.HashRing //一致性哈希
}

// HostInfoInterface interface for host
type HostInfoInterface interface {
	GetTarget(req *http.Request) string
}

//getIPAddr Get IP address
func getIPAddr(req *http.Request) []string {
	return []string{strings.Split(req.RemoteAddr, ":")[0]}
}
func RandInt() int {
	return rand.Intn(2) + 1
}

// GetTarget 选取 目标
func (hostInfo *HostInfo) GetTarget(req *http.Request) string {
	var route string
	if hostInfo.IsMultiTarget {
		if hostInfo.MultiTargetMode == SelectModeRandom { //随机模式
			route = hostInfo.MultiTarget[RandInt()%hostInfo.Length]
		} else if hostInfo.MultiTargetMode == SelectModePoll { //轮询模式
			route = hostInfo.MultiTarget[hostInfo.PoolModeIndex]
			//todo 并发问题会计算不准确，并不是真的轮训，但是应该不影响整体平衡
			//但是又不想用并发包以某种阻塞的方式去轮训
			hostInfo.PoolModeIndex++
			hostInfo.PoolModeIndex = hostInfo.PoolModeIndex % hostInfo.Length
		} else if hostInfo.MultiTargetMode == SelectModeHash { //哈希模式
			ips := getIPAddr(req)
			route, _ = hostInfo.hashRing.GetNode(ips[0])
		} else { //未配置或配置错误使用随机模式
			route = hostInfo.MultiTarget[rand.Int()%hostInfo.Length]
		}
	} else {
		route = hostInfo.Target
	}
	return route
}

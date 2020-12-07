package loadbalance

import (
	"math/rand"
	"net/http"

	"github.com/g4zhuj/hashring"
)

// 多种选择方式
const (
	SelectModeRandom int = 1 //随机选择
	SelectModePoll   int = 2 //轮询选择
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
		} else { //未配置或配置错误使用随机模式
			route = hostInfo.MultiTarget[rand.Int()%len(hostInfo.MultiTarget)]
		}
	} else {
		route = hostInfo.Target
	}
	return route
}

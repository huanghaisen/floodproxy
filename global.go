package main

import (
	"floodproxy/pkg/loadbalance"
	"github.com/spf13/viper"
)

var (
	// FloodDataConfig :
	FloodDataConfig *Server
	// FloodDataViper ：
	FloodDataViper *viper.Viper
	HostInfo       *loadbalance.HostInfo
)

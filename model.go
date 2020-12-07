package main

type Server struct {
	System    System      `json:"system" yaml:"system"`
	HttpProxy []HttpProxy `json:"httpProxy" yaml:"httpProxy"`
}
type System struct {
	AppName  string `json:"appName" yaml:"appName"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	LogLevel string `json:"log_level" yaml:"LogLevel"`
}
type HttpProxy struct {
	Target    string `json:"target" yaml:"target"`
	Proxypass string `json:"proxypass" yaml:"proxypass"`
}

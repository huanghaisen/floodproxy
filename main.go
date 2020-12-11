package main

import (
	"floodproxy/pkg/loadbalance"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const defaultConfigFile = "resource/config.yaml"

func main() {
	InitConfig(defaultConfigFile)
	if FloodDataConfig.System.Startby == 1 {
		log.Print("starting flood server")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Push signals into channel
		go run()
		sig := <-sigs
		log.Printf("Shutting down... %+v", sig)
	} else {
		log.Print("starting proxy server")
		//运行服务
		srv := new(AcceptSerevr)
		srv.runProxy(strconv.Itoa(FloodDataConfig.System.Port))
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-sigs
		log.Print("Shutting down... %v", sig)
	}
}
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}
func run() {
	go func() {
		//mux := http.NewServeMux()
		http.HandleFunc("/", sayHello)
		//http.HandleFunc("/hello", accessControl(mux))
		listen := fmt.Sprintf("%s:%d", "0.0.0.0", FloodDataConfig.System.Port)
		fmt.Println("HTTP server run on ", listen)
		if err := http.ListenAndServe(listen, nil); err != nil {
			log.Println(err.Error())
		}
		log.Println("end")
	}()
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	//n, err := fmt.Fprintln(w, "hello world")
	var result = "hello world"
	log.Print(result)
	_, _ = w.Write([]byte(result))
}

func InitConfig(configUrl string) {
	fmt.Println("==================== begin initialise config ====================")
	v := viper.New()
	v.SetConfigFile(configUrl)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&FloodDataConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&FloodDataConfig); err != nil {
		fmt.Println(err)
	}
	FloodDataViper = v
	fmt.Println("FloodDataViper:", FloodDataViper)
	fmt.Println("FloodDataConfig:", FloodDataConfig)
	HostInfo = LoadHost()
	fmt.Println("==================== end initialise config ====================")
}

// LoadHost Load host configuration
func LoadHost() *loadbalance.HostInfo {
	hostInfo := new(loadbalance.HostInfo)
	hostInfo.IsMultiTarget = true
	for _, proxy := range FloodDataConfig.HttpProxy {
		hostInfo.MultiTarget = append(hostInfo.MultiTarget, proxy.Proxypass)
	}
	hostInfo.MultiTargetMode = 1
	hostInfo.Length = len(hostInfo.MultiTarget)
	return hostInfo
}

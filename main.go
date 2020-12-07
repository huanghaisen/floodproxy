package main

import (
	"floodproxy/pkg/loadbalance"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const defaultConfigFile = "resource/config.yaml"

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Print("starting floodproxy")
	sigs := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Push signals into channel

	go run(stop, wg)

	sig := <-sigs
	log.Printf("Shutting down... %+v", sig)
	close(stop) // Tell goroutines to stop themselves
	wg.Wait()
	// InitConfig(defaultConfigFile)

	// //运行服务
	// srv := new(AcceptSerevr)
	// srv.runProxy("8009")

	// sigs := make(chan os.Signal)
	// signal.Notify(sigs, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// sig := <-sigs
	// log.Print("Shutting down... %v", sig)
}
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}
func run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		//mux := http.NewServeMux()
		http.HandleFunc("/", sayHello)
		//http.HandleFunc("/hello", accessControl(mux))
		listen := fmt.Sprintf("%s:%d", "0.0.0.0", 8006)
		fmt.Println("HTTP server run on ", listen)
		if err := http.ListenAndServe(listen, nil); err != nil {
			log.Println(err.Error())
		}
		log.Println("end")
	}()
	<-stopCh
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	//n, err := fmt.Fprintln(w, "hello world")
	_, _ = w.Write([]byte("hello world"))
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
	//	hostInfo := LoadHost()
	fmt.Println("==================== end initialise config ====================")
}

// LoadHost Load host configuration
func LoadHost() *loadbalance.HostInfo {
	hostInfo := new(loadbalance.HostInfo)
	hostInfo.IsMultiTarget = true
	for _, proxy := range FloodDataConfig.HttpProxy {
		hostInfo.MultiTarget = append(hostInfo.MultiTarget, proxy.Proxypass)
	}
	hostInfo.MultiTargetMode = 2
	return hostInfo
}

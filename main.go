package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

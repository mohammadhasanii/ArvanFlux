package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
)

type LoadBalancer struct {
	servers []string
	index   uint64
	mutex   sync.Mutex
}

func initializeLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		index:   0,
	}
}

func (lb *LoadBalancer) getNextServer() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.servers) == 0 {
		return ""
	}

	server := lb.servers[atomic.AddUint64(&lb.index, 1)%uint64(len(lb.servers))]
	return server
}

func (lb *LoadBalancer) processRequest(w http.ResponseWriter, r *http.Request) {
	server := lb.getNextServer()
	if server == "" {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	resp, err := http.Get(server + r.RequestURI)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Error copying response", http.StatusInternalServerError)
	}
}

func main() {
	servers := []string{
		"https://dejban.arvancloud.ir", // Replace with your backend servers or frontend server
		"https://dejban2.arvancloud.ir",
	}

	lb := initializeLoadBalancer(servers)
	http.HandleFunc("/", lb.processRequest)

	fmt.Println("ArvanFlux is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

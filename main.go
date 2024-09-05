package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type ServerStatus struct {
	url       string
	available bool
	lastCheck time.Time
}

type LoadBalancer struct {
	servers []ServerStatus
	index   uint64
	mutex   sync.RWMutex
}

func initializeLoadBalancer(serverURLs []string) *LoadBalancer {
	servers := make([]ServerStatus, len(serverURLs))
	for i, url := range serverURLs {
		servers[i] = ServerStatus{url: url, available: true, lastCheck: time.Now()}
	}
	return &LoadBalancer{
		servers: servers,
		index:   0,
	}
}

func (lb *LoadBalancer) getNextServer() string {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	count := uint64(len(lb.servers))
	for i := uint64(0); i < count; i++ {
		idx := (atomic.AddUint64(&lb.index, 1) - 1) % count
		server := lb.servers[idx]
		if server.available {
			return server.url
		}
	}
	return ""
}

func (lb *LoadBalancer) temporarilyDisableServer(serverURL string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i, server := range lb.servers {
		if server.url == serverURL {
			lb.servers[i].available = false
			lb.servers[i].lastCheck = time.Now()
			break
		}
	}

	go func() { //If one of the servers goes down, all requests will be transferred to the available server for 1 minute
		time.Sleep(time.Minute) //you can change it to 2 minutes    (2 * time.Minute)
		lb.mutex.Lock()
		defer lb.mutex.Unlock()
		for i, server := range lb.servers {
			if server.url == serverURL && !server.available && time.Since(server.lastCheck) >= time.Minute {
				lb.servers[i].available = true
			}
		}
	}()
}

func (lb *LoadBalancer) processRequest(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextServer()
	if targetServer == "" {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	target, err := url.Parse(targetServer)
	if err != nil {
		http.Error(w, "Error parsing target server URL", http.StatusInternalServerError)
		return
	}

	proxyReq := new(http.Request)
	*proxyReq = *r
	proxyReq.URL.Scheme = target.Scheme
	proxyReq.URL.Host = target.Host
	proxyReq.Header.Set("X-Forwarded-Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", r.URL.Scheme)

	resp, err := http.DefaultTransport.RoundTrip(proxyReq)
	if err != nil {
		lb.temporarilyDisableServer(targetServer)
		http.Error(w, "Error forwarding request to target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	//Copy the response headers to the client
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("X-Proxy-Server", "ArvanFlux 1.1.0")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func main() {
	servers := []string{
		"http://localhost:3000",
		"http://localhost:4000",
		"http://localhost:5000",
		"http://localhost:6000",
	}

	lb := initializeLoadBalancer(servers)
	http.HandleFunc("/", lb.processRequest)

	fmt.Println("ArvanFlux 1.1.0 is running on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting proxy server:", err)
	}
}

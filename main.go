package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
)

func main() {
	http.HandleFunc("/api/v0/healthchecks/liveness", LivenessHandler)
	http.HandleFunc("/api/v0/ip", WhatsMyPublicIP)
	httpErr := http.ListenAndServe(":8080", nil)
	if httpErr != nil {
		log.Fatalf("Couldn't start HTTP Server, got %v", httpErr)
	}
}

// Custom DNS Dialer to use OpenDNS
func openDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", "208.67.222.222:53")
}

// WhatsMyPublicIP attempts to connect to an external host to determine public ipv4 addr
func WhatsMyPublicIP(w http.ResponseWriter, r *http.Request) {
	resolver := net.Resolver{PreferGo: true, Dial: openDNSDialer}
	ctx := context.Background()
	ipaddr, err := resolver.LookupIPAddr(ctx, "myip.opendns.com")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to resolve addr")
	}
	reg, _ := regexp.Compile("{(.*)}")
	ipaddrStr := reg.FindStringSubmatch(fmt.Sprintf("%v", ipaddr[0]))[1]
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("%v", ipaddrStr))
}

// LivenessHandler prints OK with status HTTP 200 to show the application is running
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `OK`)
}

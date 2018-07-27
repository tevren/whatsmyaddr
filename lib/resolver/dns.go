package resolver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
)

// Custom DNS Dialer to use OpenDNS
func openDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, network, "208.67.222.222:53")
	if err != nil {
		return nil, errors.New("Couldn't create dial context")
	}
	return conn, nil
}

func GetPublicAddrWithDNS() (string, error) {
	resolver := net.Resolver{PreferGo: true, Dial: openDNSDialer}
	ctx := context.Background()
	ipaddr, _ := resolver.LookupIPAddr(ctx, "myip.opendns.com")
	reg, _ := regexp.Compile("{(.*)}")
	ipaddrStr := reg.FindStringSubmatch(fmt.Sprintf("%v", ipaddr[0]))[1]
	return ipaddrStr, nil
}

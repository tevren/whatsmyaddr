package resolver

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetPublicAddrWithHTTP() (string, error) {
	resp, _ := http.Get("http://icanhazip.com")
	body, _ := ioutil.ReadAll(resp.Body)
	ipaddr := strings.TrimSuffix(string(body), "\n")
	if net.ParseIP(ipaddr) == nil {
		return "", errors.New("Couldn't resolve public ip address")
	}
	return ipaddr, nil
}

package youtube

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Proxy struct {
	ID       int
	Login    string
	Password string
	IP       string
	Port     string
}

type ProxyOption struct {
	Path string
	Key  int
}

func (proxy Proxy) GetProxyString() string {
	return fmt.Sprintf("%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.IP, proxy.Port)
}

func ParsFile(pathFile string) []Proxy {

	proxyList := make([]Proxy, 0)
	proxyFile, _ := ioutil.ReadFile(pathFile)

	proxyLines := strings.Split(string(proxyFile), "\n")
	for i := 0; i < len(proxyLines); i++ {
		if proxyLines[i] != "" {
			proxyLine := strings.Split(proxyLines[i], ":")
			newProxy := Proxy{ID: i, Login: proxyLine[0], Password: proxyLine[1], IP: proxyLine[2], Port: proxyLine[3]}
			proxyList = append(proxyList, newProxy)
		}
	}

	return proxyList
}

func GetProxy(option *ProxyOption) Proxy {
	collection := ParsFile(option.Path)

	if len(collection)-1 == option.Key {

		proxy := collection[option.Key]
		option.Key = 0

		return proxy
	} else {
		proxy := collection[option.Key]
		option.Key++

		return proxy
	}
}

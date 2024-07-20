package xnet

import (
	"net/url"
	"strconv"
	"strings"
)

// 从一个url中提取出host与port数据,支持http(s)://host:port,host:port,host等格式
// 如http://192.168.91.3:83;https://p.abmp.cc;192.168.91.3:83;192.168.91.3
func ParseHostAndPort(urlHost string) (host []byte, port *int) {
	if len(urlHost) <= 0 {
		return host, nil
	}
	urlHost = strings.ToLower(urlHost)
	urlHost = strings.TrimPrefix(urlHost, "http://")
	urlHost = strings.TrimPrefix(urlHost, "https://")
	urlSegementList := strings.Split(urlHost, ":")
	if len(urlSegementList) <= 0 {
		return nil, nil
	}
	host = []byte(urlSegementList[0])
	if len(urlSegementList) > 1 {
		portValue, err := strconv.Atoi(urlSegementList[1])
		if err != nil {
			return host, nil
		}
		port = &portValue
	}
	return host, port
}

// 更改url的hostname
func ChangeUrlHost(urlValue string, newHostName string) string {
	url, err := url.Parse(urlValue)
	if err != nil {
		return urlValue
	}
	_, port := ParseHostAndPort(urlValue)
	url.Host = string(newHostName) + ":" + strconv.Itoa(*port)
	return url.String()
}

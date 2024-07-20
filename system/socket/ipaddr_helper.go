package socket

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
)

func FormatAddressPort(address string, port int) string {
	return net.JoinHostPort(address, strconv.Itoa(port))
}

// 检测指定的ip地址是否是IPv4或者IPv6 any地址
// ip只能是*net.IP或者string类型，其它类型会抛出panic
func IsAny(ip interface{}) bool {
	return IsAnyV4(ip) || IsAnyV6(ip)
}

// 检测指定的ip地址是否是IPv4 any地址
// ip只能是*net.IP或者string类型，其它类型会抛出panic
func IsAnyV4(ip interface{}) bool {
	return iptos(ip) == "0.0.0.0"
}

func IsAnyV6(ip interface{}) bool {
	ips := iptos(ip)
	return ips == "::" || ips == "[::]"
}

func iptos(ip interface{}) string {
	if ip == nil || reflect.TypeOf(ip).Kind() == reflect.Ptr && reflect.ValueOf(ip).IsNil() {
		return ""
	}
	switch x := ip.(type) {
	case string:
		return x
	case *string:
		if x == nil {
			return ""
		}
		return *x
	case net.IP:
		return x.String()
	case *net.IP:
		return x.String()
	case *net.IPAddr:
		return x.IP.String()
	case *net.TCPAddr:
		return x.IP.String()
	case *net.UDPAddr:
		return x.IP.String()
	default:
		panic(fmt.Sprintf("invalid type: %T", ip))
	}
}

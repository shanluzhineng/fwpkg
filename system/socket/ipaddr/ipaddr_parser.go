package ipaddr

import (
	"fmt"
	"net"
	"strings"

	"github.com/shanluzhineng/fwpkg/system/multierror"
	"github.com/shanluzhineng/fwpkg/system/socket/template"
)

type IPAddrTemplate struct {
	//保存执行过程中的错误信息
	Err error
}

// 根据模板来解析ip地址
// 返回一个 *net.IPAddr and *net.UnixAddr对象数据
// 如果有错误，会设置t.Err属性
func (t *IPAddrTemplate) ExpandAddrs(name string, s *string) []net.Addr {
	t.Err = nil
	if s == nil || *s == "" {
		return nil
	}

	x, err := template.Parse(*s)
	if err != nil {
		t.Err = multierror.Append(t.Err, fmt.Errorf("%s: error parsing %q: %s", name, *s, err))
		return nil
	}

	var addrs []net.Addr
	for _, a := range strings.Fields(x) {
		switch {
		case strings.HasPrefix(a, "unix://"):
			addrs = append(addrs, &net.UnixAddr{Name: a[len("unix://"):], Net: "unix"})
		default:
			// net.ParseIP does not like '[::]'
			ip := net.ParseIP(a)
			if a == "[::]" {
				ip = net.ParseIP("::")
			}
			if ip == nil {
				t.Err = multierror.Append(t.Err, fmt.Errorf("%s: invalid ip address: %s", name, a))
				return nil
			}
			addrs = append(addrs, &net.IPAddr{IP: ip})
		}
	}

	return addrs
}

// 取第一个地址，如果机机存在着多个地址，则会设置t.Err的错误并返回nil
func (t *IPAddrTemplate) ExpandFirstAddr(name string, s *string) net.Addr {
	t.Err = nil
	if s == nil || *s == "" {
		return nil
	}

	addrs := t.ExpandAddrs(name, s)
	if len(addrs) == 0 {
		return nil
	}
	if len(addrs) > 1 {
		var x []string
		for _, a := range addrs {
			x = append(x, a.String())
		}
		t.Err = multierror.Append(t.Err, fmt.Errorf("%s: 发现多个ip地址: %s", name, strings.Join(x, " ")))
		return nil
	}
	return addrs[0]
}

// 取第一个ip地址，如果机机存在着多个地址，则会设置t.Err的错误并返回nil
func (t *IPAddrTemplate) ExpandFirstIP(name string, s *string) *net.IPAddr {
	t.Err = nil
	if s == nil || *s == "" {
		return nil
	}

	addr := t.ExpandFirstAddr(name, s)
	if addr == nil {
		return nil
	}
	switch a := addr.(type) {
	case *net.IPAddr:
		return a
	case *net.UnixAddr:
		t.Err = multierror.Append(t.Err, fmt.Errorf("%s cannot be a unix socket", name))
		return nil
	default:
		t.Err = multierror.Append(t.Err, fmt.Errorf("%s has invalid address type %T", name, a))
		return nil
	}
}

// 如果pri不为空，则返回pri，否则返回sec参数
func ChooseNotNilIPAddr(pri *net.IPAddr, sec *net.IPAddr) *net.IPAddr {
	if pri != nil {
		return pri
	}
	return sec
}

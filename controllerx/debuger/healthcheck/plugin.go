//go:build !windows
// +build !windows

package healthcheck

import (
	"fmt"
)

func init() {
	fmt.Printf("plugin healthcheck init function called\r\n")
}

type Bootstrap struct {
}

func newBootstrap() Bootstrap {
	b := Bootstrap{}
	return b
}

func (b Bootstrap) BootstrapPlugin() (err error) {
	return nil
}

var PluginBootstrap = newBootstrap()

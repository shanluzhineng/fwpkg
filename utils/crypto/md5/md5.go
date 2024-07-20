// Package md5 provides md5 encryption utilities
package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Encrypt md5 encryption util
func Encrypt(in string) (out string) {
	md5Ctx := md5.New()                // md5 init
	n, err := md5Ctx.Write([]byte(in)) // md5 update
	if err == nil && n != 0 {
		cipherStr := md5Ctx.Sum(nil)        // md5 final
		out = hex.EncodeToString(cipherStr) // hex digest
	}
	return
}

package base64

import (
	"testing"

	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestBase64String(t *testing.T) {
	src := "Hello world"
	str := EncodeToString(src)
	decoded, err := DecodeToString(str)
	assert.Equal(t, nil, err)
	log.Debugf("src: %v, base64: %v, decoded: %v", src, str, decoded)
	assert.Equal(t, src, decoded)
}

func TestBase64Bytes(t *testing.T) {
	src := []byte("Hello world")
	str := Encode(src)
	decoded, err := Decode(str)
	assert.Equal(t, nil, err)
	log.Debugf("src: %v, base64: %v, decoded: %v", src, str, decoded)
	assert.Equal(t, src, decoded)
}

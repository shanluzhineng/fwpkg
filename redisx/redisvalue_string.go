package redisx

import (
	"strconv"
	"time"
)

type redisStringValue struct {
	data []byte
	err  error

	unmarshal UnmarshalFunc
}

var _ IRedisValue = (*redisStringValue)(nil)

func newRedisValue(data []byte, unmarshal UnmarshalFunc) *redisStringValue {
	return &redisStringValue{
		data:      data,
		unmarshal: unmarshal,
	}
}

// 构建IRedisValue的nullable模式实现
func newNilRedisValue() *redisStringValue {
	return &redisStringValue{
		data: nil,
		err:  nil,
	}
}

func newErrRedisValue(err error) *redisStringValue {
	return &redisStringValue{
		data: nil,
		err:  err,
	}
}

func (v *redisStringValue) Exist() bool {
	return v.data != nil
}

func (v *redisStringValue) Bytes() []byte {
	return v.data
}

func (v *redisStringValue) Err() error {
	return v.err
}

func (v *redisStringValue) ToValue(val interface{}) error {
	if v.err != nil {
		return v.err
	}
	if v.data == nil {
		return ErrValueIsNil
	}
	return v.unmarshal(v.data, val)
}

func (v *redisStringValue) ValToString() string {
	if v.data == nil {
		return ""
	}
	return string(v.data)
}

func (v *redisStringValue) ValToInt() (int, error) {
	if v.err != nil {
		return 0, v.err
	}
	s := v.ValToString()
	if len(s) <= 0 {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func (v *redisStringValue) ValToInt32() (int32, error) {
	value, err := v.ValToInt()
	if err != nil {
		return 0, err
	}
	return int32(value), err
}

func (v *redisStringValue) ValToInt64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	s := v.ValToString()
	if len(s) <= 0 {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func (v *redisStringValue) ValToBool() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	s := v.ValToString()
	if len(s) <= 0 {
		return false, nil
	}
	return strconv.ParseBool(s)
}

func (v *redisStringValue) ValToTime() (time.Time, error) {
	if v.err != nil {
		return time.Time{}, v.err
	}
	s := v.ValToString()
	if len(s) <= 0 {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, s)
}

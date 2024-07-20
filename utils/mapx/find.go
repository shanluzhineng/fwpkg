package mapx

import "strings"

func KeyExists(k string, m map[string]interface{}, keyInsensitivise bool) string {
	lk := k
	if keyInsensitivise {
		lk = strings.ToLower(k)
	}
	for mk := range m {
		lmk := mk
		if keyInsensitivise {
			lmk = strings.ToLower(mk)
		}
		if lmk == lk {
			return mk
		}
	}
	return ""
}

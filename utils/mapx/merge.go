package mapx

import "github.com/shanluzhineng/fwpkg/utils/stringslice"

type MergeConfig struct {
	OnlyReplaceExist bool
	OnlyAdd          bool
	KeyInsensitivise bool
	ExcludeKey       []string
}

func defaultMergeConfig() *MergeConfig {
	return &MergeConfig{
		OnlyReplaceExist: false,
		OnlyAdd:          false,
		KeyInsensitivise: false,
		ExcludeKey:       make([]string, 0),
	}
}

// merge map
func MergeMaps(src map[string]interface{}, dst map[string]interface{}, opts ...MergeConfig) {
	currentOpt := defaultMergeConfig()
	if len(opts) > 0 {
		currentOpt = &opts[0]
	}
	for sk, sv := range src {
		if len(currentOpt.ExcludeKey) > 0 {
			ignore := false
			if currentOpt.KeyInsensitivise {
				ignore = stringslice.ContainsIgnoreLowercase(currentOpt.ExcludeKey, sk)
			} else {
				ignore = stringslice.Contains(currentOpt.ExcludeKey, sk)
			}
			if ignore {
				continue
			}
		}
		tk := KeyExists(sk, dst, currentOpt.KeyInsensitivise)
		if tk == "" {
			// not exist
			if currentOpt.OnlyReplaceExist {
				continue
			}
			dst[sk] = sv
			continue
		}
		if currentOpt.OnlyAdd {
			continue
		}
		dst[tk] = sv
	}
}

func MergeMapsT[TS comparable, VS any, TD comparable, VD any](srcMap map[TS]VS, toMap map[TD]VD, itemFunc func(srcKey TS, srcValue VS) (TD, VD)) {
	for sk, sv := range srcMap {
		toKey, toValue := itemFunc(sk, sv)
		toMap[toKey] = toValue
	}
}

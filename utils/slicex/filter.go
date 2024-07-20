package slicex

// filter slice
func Filter[SV any](list []SV, filter func(item SV) bool) []SV {
	result := make([]SV, 0)
	for _, eachSV := range list {
		inFilter := filter(eachSV)
		if !inFilter {
			continue
		}
		result = append(result, eachSV)
	}
	return result
}

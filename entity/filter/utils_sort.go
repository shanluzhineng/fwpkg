package filter

import (
	"encoding/json"

	"github.com/shanluzhineng/fwpkg/entity"
	"go.mongodb.org/mongo-driver/bson"
)

// GetSorts Get entity.Sort
func GetSorts(getKeyFn func(key string) string) (sorts []entity.Sort, err error) {
	// bind
	sortStr := getKeyFn(SortQueryField)
	if err := json.Unmarshal([]byte(sortStr), &sorts); err != nil {
		return nil, err
	}
	return sorts, nil
}

// GetSortsOption Get entity.Sort
func GetSortsOption(getKeyFn func(key string) string) (sort bson.D, err error) {
	sorts, err := GetSorts(getKeyFn)
	if err != nil {
		return nil, err
	}

	// if len(sorts) == 0 {
	// 	return bson.D{{Key: "_id", Value: -1}}, nil
	// }

	return SortsToOption(sorts)
}

func MustGetSortOption(getKeyFn func(key string) string) (sort bson.D) {
	sort, err := GetSortsOption(getKeyFn)
	if err != nil {
		return nil
	}
	return sort
}

// SortsToOption Translate entity.Sort to bson.D
func SortsToOption(sorts []entity.Sort) (sort bson.D, err error) {
	sort = bson.D{}
	for _, s := range sorts {
		switch s.Direction {
		case ASCENDING:
			sort = append(sort, bson.E{Key: s.Key, Value: 1})
		case DESCENDING:
			sort = append(sort, bson.E{Key: s.Key, Value: -1})
		}
	}
	// if len(sort) == 0 {
	// 	sort = bson.D{{Key: "_id", Value: -1}}
	// }
	return sort, nil
}

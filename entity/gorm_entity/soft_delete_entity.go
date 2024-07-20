package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// 支持软删除的实体
type SoftDeleteEntity struct {
	//是否已删除
	IsDeleted bool `json:"isDeleted"`
	//删除人员
	DeleterId *uuid.UUID `json:"deleterId"`
	//删除时间
	DeletionTime time.Time `json:"deletionTime"`
}

package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// 可审核的entity
type AuditedEntity struct {
	Entity
	//创建时间
	CreationTime time.Time `json:"creationTime" bson:"creationTime" gorm:"column:CreationTime"`
	//创建人
	CreatorId *uuid.UUID `json:"creatorId" bson:"creatorId" gorm:"column:CreatorId"`
	//最后修改时间
	LastModificationTime time.Time `json:"lastModificationTime"`
	//最后修改人id
	LastModifierId *uuid.UUID `json:"lastModifierId"`
}

// 创建时设置对象的基本信息
func (entity *AuditedEntity) BeforeCreate() {
	//调用基类的函数
	entity.Entity.BeforeCreate()
	if entity.CreationTime.IsZero() {
		entity.CreationTime = time.Now()
	}
}

func (entity *AuditedEntity) BeforeUpdate() {
	//设置最后修改时间
	entity.LastModificationTime = time.Now()
}

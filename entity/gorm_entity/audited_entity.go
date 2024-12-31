package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
func (entity *AuditedEntity) BeforeCreate(db *gorm.DB) error {
	//调用基类的函数
	entity.Entity.BeforeCreate(db)
	if entity.CreationTime.IsZero() {
		entity.CreationTime = time.Now()
	}
	return nil
}

func (entity *AuditedEntity) BeforeUpdate(db *gorm.DB) error {
	//设置最后修改时间
	entity.LastModificationTime = time.Now()
	return nil
}

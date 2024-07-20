package entity

import (
	uuid "github.com/satori/go.uuid"
)

//一个抽象的entity对象
type Entity struct {
	Id uuid.UUID `json:"id"`
}

//自动创建id
func (entity *Entity) BeforeCreate() {
	if entity.Id != uuid.Nil {
		//已经创建了，则不再创建
		return
	}
	//创建一个新的id
	entity.Id = uuid.NewV4()
}

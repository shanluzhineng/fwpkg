package controllerx

import (
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/entity"
	"github.com/shanluzhineng/fwpkg/mongodbr"
)

func GetEntityService[T mongodbr.IEntity]() entity.IEntityService[T] {
	return app.Context.GetInstance(new(entity.IEntityService[T])).(entity.IEntityService[T])
}

package controllerx

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/shanluzhineng/fwpkg/controllerx/fwauth"
	"github.com/shanluzhineng/fwpkg/controllerx/responsex"
	"github.com/shanluzhineng/fwpkg/entity"
	"github.com/shanluzhineng/fwpkg/entity/filter"
	"github.com/shanluzhineng/fwpkg/mongodbr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EntityController[T mongodbr.IEntity] struct {
	RouterPath    string
	EntityService entity.IEntityService[T]

	Options BaseEntityControllerOptions
	once    sync.Once
}

func (c *EntityController[T]) RegistRouter(webapp *IrisApplication, opts ...BaseEntityControllerOption) router.Party {
	for _, eachOpt := range opts {
		eachOpt(&(c.Options))
	}

	routerParty := webapp.Party(c.RouterPath)

	if !c.Options.AllDisabled {
		routerParty.Get("/all", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.All)...)
	}
	if !c.Options.ListDisabled {
		routerParty.Get("/", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.GetList)...)
	}
	if !c.Options.GetByIdDisabled {
		routerParty.Get("/{id}", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.GetById)...)
	}
	if !c.Options.CreateDisabled {
		routerParty.Post("/", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.Create)...)
	}
	if !c.Options.UpdateDisabled {
		routerParty.Put("/{id}", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.Update)...)
	}
	if !c.Options.DeleteDisabled {
		routerParty.Delete("/{id}", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.Delete)...)
	}
	if !c.Options.DeleteListDisabled {
		routerParty.Delete("/", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.DeleteList)...)
	}

	return routerParty
}

func (c *EntityController[T]) MergeAuthenticatedContextIfNeed(authenticatedDisabled bool, handlers ...context.Handler) []context.Handler {
	handlerList := make([]context.Handler, 0)
	if !authenticatedDisabled {
		// handler auth
		handlerList = append(handlerList, fwauth.GetCasdoorMiddleware().Serve)
	}
	handlerList = append(handlerList, handlers...)
	return handlerList
}

func (c *EntityController[T]) GetEntityService() entity.IEntityService[T] {
	c.once.Do(func() {
		if c.EntityService != nil {
			return
		}
		c.EntityService = GetEntityService[T]()
	})
	return c.EntityService
}

func (c *EntityController[T]) All(ctx iris.Context) {
	filter := map[string]interface{}{}
	if !c.Options.FilterCurrentUserForListDisabled {
		// auto filter current userId
		AddUserIdFilterIfNeed(filter, new(T), ctx)
	}

	if c.Options.ListFilterFunc != nil {
		c.Options.ListFilterFunc(new(T), filter, ctx)
	}
	var list []T
	var err error
	if len(filter) > 0 {
		list, err = c.GetEntityService().FindList(filter)
	} else {
		list, err = c.GetEntityService().FindAll()
	}
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccessWithListData(ctx, list, int64(len(list)))
}

func (c *EntityController[T]) GetList(ctx iris.Context) {
	all := filter.MustGetFilterAll(ctx.FormValue)
	if all {
		c.All(ctx)
		return
	}

	// params
	pagination := MustGetPagination(ctx)
	query := filter.MustGetFilterQuery(ctx.FormValue)
	sort := filter.MustGetSortOption(ctx.FormValue)

	if !c.Options.FilterCurrentUserForListDisabled {
		// auto filter current userId
		AddUserIdFilterIfNeed(query, new(T), ctx)
	}
	service := c.GetEntityService()
	list, err := service.FindList(query, mongodbr.FindOptionWithSort(sort),
		mongodbr.FindOptionWithPage(int64(pagination.Page), int64(pagination.Size)))
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}

	count, err := service.Count(query)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccessWithListData(ctx, list, count)
}

// get by id
func (c *EntityController[T]) GetById(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		responsex.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}

	id, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id,id must be bson id format,id:%s", idValue))
		return
	}
	item, err := c.GetEntityService().FindById(id)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		responsex.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
		return
	}
	// // filter user is current user
	// if !FilterMustIsCurrentUserId(item, ctx) {
	// 	responsex.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
	// 	return
	// }
	responsex.HandleSuccessWithData(ctx, item)
}

// create
func (c *EntityController[T]) Create(ctx iris.Context) {
	input := new(T)
	err := ctx.ReadJSON(&input)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, err)
		return
	}
	err = mongodbr.Validate(input)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, err)
		return
	}

	// handler user info
	c.SetUserInfo(ctx, input)

	newItem, err := c.GetEntityService().Create(input)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccessWithData(ctx, newItem)
}

// update
func (c *EntityController[T]) Update(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		responsex.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}
	id, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id,id must be bson id format,id:%s", idValue))
		return
	}
	service := c.GetEntityService()
	item, err := service.FindById(id)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		responsex.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}
	// filter user is current user
	// if !FilterMustIsCurrentUserId(item, ctx) {
	// 	responsex.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
	// 	return
	// }

	input := make(map[string]interface{})
	err = ctx.ReadJSON(&input)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, err)
		return
	}

	err = service.UpdateFields(id, input)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccess(ctx)
}

// delete
func (c *EntityController[T]) Delete(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		responsex.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}
	oid, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id format,err:%s", err.Error()))
		return
	}
	service := c.GetEntityService()
	item, err := service.FindById(oid)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	if item == nil {
		responsex.HandleErrorBadRequest(ctx, fmt.Errorf("not found item,id:%s", idValue))
		return
	}
	// filter user is current user
	// if !FilterMustIsCurrentUserId(item, ctx) {
	// 	responsex.HandleErrorInternalServerError(ctx, fmt.Errorf("invalid id,id:%s", idValue))
	// 	return
	// }

	err = c.GetEntityService().Delete(oid)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccess(ctx)
}

// delete
func (c *EntityController[T]) DeleteList(ctx iris.Context) {
	payload, err := GetBatchRequestPayload(ctx)
	if err != nil {
		responsex.HandleErrorBadRequest(ctx, err)
		return
	}
	if len(payload.Ids) <= 0 {
		responsex.HandleSuccess(ctx)
		return
	}
	filter := bson.M{
		"_id": bson.M{"$in": payload.Ids},
	}
	// auto filter current userId
	// AddUserIdFilterIfNeed(filter, new(T), ctx)

	_, err = c.GetEntityService().DeleteMany(filter)
	if err != nil {
		responsex.HandleErrorInternalServerError(ctx, err)
		return
	}
	responsex.HandleSuccess(ctx)
}

func (c *EntityController[T]) SetUserInfo(ctx iris.Context, entityValue interface{}) {
	userinfoProvider, ok := entityValue.(entity.IEntityWithUser)
	if !ok {
		return
	}
	userId := GetUserId(ctx)
	if userId != "" {
		userinfoProvider.SetUserCreator(userId)
	}
}

package controllerx

import (
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	"github.com/shanluzhineng/fwpkg/controllerx/fwauth"
	"github.com/shanluzhineng/fwpkg/entity"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/system/reflector"
)

// user Custom claims struct
type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	TenantId    uuid.UUID // 租户id，一般为companyId，uniweb中新增用户与casdoorId相同
	NickName    string
	AuthorityId string
	CasClaim    *casdoorsdk.Claims
	Role        string // admin / normal
}
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

// 请求上下文数据
type BaseServiceContext struct {
	//当前用户信息
	CurrentUserInfo CustomClaims
}

// Gin operations
func GetClaims(c *gin.Context) (*CustomClaims, error) {
	var cusClaim *CustomClaims
	if claimAny, ok := c.Get(fwauth.AuthKey); ok {
		claim := claimAny.(*casdoorsdk.Claims)
		cusClaim = &CustomClaims{}
		cusClaim.CasClaim = claim
		cusClaim.Id = claim.Id
		cusClaim.Username = claim.DisplayName
	}
	return cusClaim, nil
}

func SetBaseServiceContext(context *BaseServiceContext, c *gin.Context) {
	claims, err := GetClaims(c)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("[SetBaseServiceContext]>> GetClaims failed: %s", err.Error()))
		return
	}
	if claims != nil {
		context.CurrentUserInfo = *claims
	}
}

// Iris operations
func GetUserId(ctx iris.Context) string {
	claims := fwauth.GetCasdoorMiddleware().GetUserClaims(ctx)
	if claims != nil {
		return claims.Id
	}
	return ""
}

func checkEntityIsIEntityWithUser(entityValue interface{}) entity.IEntityWithUser {
	v, ok := entityValue.(entity.IEntityWithUser)
	if !ok {
		return nil
	}
	return v
}

func AddUserIdFilterIfNeed(filter map[string]interface{}, entity interface{}, ctx iris.Context) {
	if filter == nil {
		return
	}
	if checkEntityIsIEntityWithUser(entity) == nil {
		return
	}
	currentUserId := GetUserId(ctx)
	if currentUserId == "" {
		return
	}
	filter["creatorId"] = currentUserId
}

func FilterMustIsCurrentUserId(entity interface{}, ctx iris.Context) bool {
	entityWithUser := checkEntityIsIEntityWithUser(entity)
	if entityWithUser == nil {
		return true
	}
	userId := GetUserId(ctx)
	if userId == "" {
		return false
	}
	ok := (userId == entityWithUser.GetCreatorId())
	if !ok {
		log.Warn(fmt.Sprintf("用户修改的对象不属于当前登录的用户,对象类型:%s,userId:%s",
			reflector.GetFullName(entity),
			userId))
	}
	return ok
}

package entity

import uuid "github.com/satori/go.uuid"

//带租户的实体
type MayHaveTenantEntity struct {
	TenantId *uuid.UUID `json:"tenantId"`
}

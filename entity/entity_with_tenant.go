package entity

type IEntityWithTenant interface {
	SetTenantId(tenantId string)
	GetTenantId() string
}

type EntityWithTenant struct {
	TenantId string `json:"tenantId" bson:"tenantId"`
}

// #region IEntityWithUser Members

func (t *EntityWithTenant) SetTenantId(tenantId string) {
	t.TenantId = tenantId
}

func (p *EntityWithTenant) GetTenantId() string {
	return p.TenantId
}

// #endregion

func CheckEntityIsIEntityWithTenant(entityValue interface{}) IEntityWithTenant {
	v, ok := entityValue.(IEntityWithTenant)
	if !ok {
		return nil
	}
	return v
}

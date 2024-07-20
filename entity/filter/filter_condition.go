package filter

type FilterCondition interface {
	GetKey() (key string)
	SetKey(key string)
	GetOp() (op string)
	SetOp(op string)
	GetValue() (value interface{})
	SetValue(value interface{})
}

type Condition struct {
	Key   string      `json:"key"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

// #region FilterCondition Members

func (c *Condition) GetKey() (key string) {
	return c.Key
}

func (c *Condition) SetKey(key string) {
	c.Key = key
}

func (c *Condition) GetOp() (op string) {
	return c.Op
}

func (c *Condition) SetOp(op string) {
	c.Op = op
}

func (c *Condition) GetValue() (value interface{}) {
	return c.Value
}

func (c *Condition) SetValue(value interface{}) {
	c.Value = value
}

// #endregion

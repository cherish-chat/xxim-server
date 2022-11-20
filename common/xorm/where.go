package xorm

type GormWhere struct {
	Where string
	args  []interface{}
}

func NewGormWhere() []GormWhere {
	return make([]GormWhere, 0)
}

func Where(where string, args ...interface{}) GormWhere {
	return GormWhere{Where: where, args: args}
}

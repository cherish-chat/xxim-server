package utils

import "go.mongodb.org/mongo-driver/bson"

type xMap struct {
}

var Map = &xMap{}

func (x *xMap) SS2SA(ss map[string]string) map[string]any {
	m := make(map[string]any)
	for k, v := range ss {
		m[k] = v
	}
	return m
}

func (x *xMap) SSMFromString(s string) SSM {
	ssm := make(SSM)
	_ = Json.Unmarshal([]byte(s), &ssm)
	return ssm
}

type SSM map[string]string

func (ssm SSM) ToSA() map[string]any {
	m := make(map[string]any)
	for k, v := range ssm {
		m[k] = v
	}
	return m
}

func (ssm SSM) ToSS() map[string]string {
	return ssm
}

func (ssm SSM) Get(k string) string {
	v, _ := ssm[k]
	return v
}

func (ssm SSM) GetInt64(k string) int64 {
	v, _ := ssm[k]
	return Number.Any2Int64(v)
}

func (ssm SSM) GetOrDefault(k string, def string) string {
	v, ok := ssm[k]
	if !ok {
		return def
	}
	return v
}

func (ssm SSM) Marshal() string {
	return Json.MarshalToString(ssm)
}

func NewSSMFromBsonM(m bson.M) SSM {
	ssm := make(SSM)
	for k, v := range m {
		ssm[k] = AnyString(v)
	}
	return ssm
}

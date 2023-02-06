package xorm

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type M map[string]any

func (m M) Get(key string, defaultValue string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return defaultValue
}

func (m M) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *M) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	if b, ok := src.([]byte); !ok {
		return nil
	} else if string(b) == "" {
		return nil
	} else if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(src.([]byte), m)
}

func (m M) StringMap() map[string]string {
	mm := make(map[string]string)
	for k, v := range m {
		mm[k] = utils.AnyToString(v)
	}
	return mm
}

type SliceString []string

func (s SliceString) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SliceString) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), s)
}

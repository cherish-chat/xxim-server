package pb

import (
	"database/sql/driver"
	"encoding/json"
)

type RequestAddFriendExtraList []*RequestAddFriendExtra

func (x RequestAddFriendExtraList) Value() (driver.Value, error) {
	return json.Marshal(x)
}

func (x *RequestAddFriendExtraList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), x)
}

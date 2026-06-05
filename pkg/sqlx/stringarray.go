package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var _ driver.Valuer = StringArray{}
var _ sql.Scanner = (*StringArray)(nil)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return a.JSON(), nil
}

func (a *StringArray) Scan(src any) error {
	switch value := src.(type) {
	case nil:
		*a = nil
		return nil
	case string:
		return a.UnmarshalJSON([]byte(value))
	case []byte:
		return a.UnmarshalJSON(value)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", src)
	}
}

func (a StringArray) JSON() string {
	if a == nil {
		return "[]"
	}
	data, _ := json.Marshal([]string(a))
	return string(data)
}

func (a StringArray) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(a))
}

func (a *StringArray) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*a = nil
		return nil
	}
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	*a = values
	return nil
}

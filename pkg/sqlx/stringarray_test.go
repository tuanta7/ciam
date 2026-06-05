package sqlx

import (
	"encoding/json"
	"testing"
)

func TestStringArrayJSON(t *testing.T) {
	values := StringArray{"openid", "email"}

	if got := values.JSON(); got != `["openid","email"]` {
		t.Fatalf("JSON() = %s", got)
	}
}

func TestStringArrayJSONNilAsEmptyArray(t *testing.T) {
	var values StringArray

	if got := values.JSON(); got != `[]` {
		t.Fatalf("JSON() = %s", got)
	}

	data, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `[]` {
		t.Fatalf("MarshalJSON() = %s", data)
	}
}

func TestStringArrayScan(t *testing.T) {
	var values StringArray

	if err := values.Scan(`["openid","email"]`); err != nil {
		t.Fatal(err)
	}

	if len(values) != 2 || values[0] != "openid" || values[1] != "email" {
		t.Fatalf("Scan() = %#v", values)
	}
}

func TestStringArrayValue(t *testing.T) {
	values := StringArray{"openid", "email"}

	got, err := values.Value()
	if err != nil {
		t.Fatal(err)
	}
	if got != `["openid","email"]` {
		t.Fatalf("Value() = %v", got)
	}
}

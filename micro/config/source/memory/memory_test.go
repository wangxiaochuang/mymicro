package memory

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
)

func TestMemory_Read(t *testing.T) {
	data := []byte(`{"database":{"host":"localhost","password":"password","datasource":"user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local"}}`)

	source := NewSource(WithJSON(data))
	c, err := source.Read()
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, c, nil)

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Fatal(err)
	}

	actualDB := actual["database"].(map[string]interface{})
	assert.Equal(t, actualDB["host"], "localhost")
}

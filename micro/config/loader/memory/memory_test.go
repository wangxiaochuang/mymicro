package memory

import (
	"testing"

	"github.com/wxc/micro/config/loader"
	sm "github.com/wxc/micro/config/source/memory"
)

func TestMy(t *testing.T) {
	data := []byte(`{"database":{"host":"localhost","password":"password","datasource":"user:password@tcp(localhost:port)/db?charset=utf8mb4&parseTime=True&loc=Local"}}`)
	m := NewLoader(loader.WithSource(sm.NewSource(sm.WithJSON(data))))
	_, err := m.Snapshot()
	if err != nil {
		t.Fatalf("get snapshot error: %s", err.Error())
	}
}

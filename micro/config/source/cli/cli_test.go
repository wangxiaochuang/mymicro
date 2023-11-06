package cli

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
	"github.com/wxc/micro/config/source"
	"github.com/wxc/micro/util/cmd"
)

func test(t *testing.T, withContext bool) {
	var src source.Source

	app := cmd.App()
	app.Name = "testapp"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "db-host",
			EnvVars: []string{"db-host"},
			Value:   "myval",
		},
	}

	if withContext {
		app.Action = func(c *cli.Context) error {
			src = WithContext(c)
			return nil
		}

		app.Run([]string{"run", "-db-host", "localhost"})
	} else {
		os.Args = []string{"run", "-db-host", "localhost"}
		src = NewSource()
	}

	c, err := src.Read()
	if err != nil {
		t.Error(err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["db"].(map[string]interface{})
	if actualDB["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", actualDB["name"])
	}
}

func TestCliSource(t *testing.T) {
	// without context
	test(t, false)
}

func TestCliSourceWithContext(t *testing.T) {
	// with context
	test(t, true)
}

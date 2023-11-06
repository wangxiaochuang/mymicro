package cli

import (
	"flag"
	"io"
	"os"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/urfave/cli/v2"
	"github.com/wxc/micro/config/source"
	"github.com/wxc/micro/util/cmd"
)

type cliSource struct {
	opts source.Options
	ctx  *cli.Context
}

func (c *cliSource) Read() (*source.ChangeSet, error) {
	var changes map[string]interface{}

	for _, f := range c.ctx.App.Flags {
		name := f.Names()[0]
		tmp := toEntry(name, c.ctx.Generic(name))
		if err := mergo.Map(&changes, tmp, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	b, err := c.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    c.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    c.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func toEntry(name string, v interface{}) map[string]interface{} {
	n := strings.ToLower(name)
	keys := strings.FieldsFunc(n, split)
	reverse(keys)
	tmp := make(map[string]interface{})
	for i, k := range keys {
		if i == 0 {
			tmp[k] = v
			continue
		}

		tmp = map[string]interface{}{k: tmp}
	}
	return tmp
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func split(r rune) bool {
	return r == '-' || r == '_'
}

func (c *cliSource) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

// Write is unsupported.
func (c *cliSource) Write(cs *source.ChangeSet) error {
	return nil
}

func (c *cliSource) String() string {
	return "cli"
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	var ctx *cli.Context

	if c, ok := options.Context.Value(contextKey{}).(*cli.Context); ok {
		ctx = c
	} else {
		app := cmd.App()
		flags := app.Flags
		set := flag.NewFlagSet(app.Name, flag.ContinueOnError)
		for _, f := range flags {
			f.Apply(set)
		}
		set.SetOutput(io.Discard)
		set.Parse(os.Args[1:])
		normalizeFlags(app.Flags, set)
		ctx = cli.NewContext(app, set, nil)
	}
	return &cliSource{
		ctx:  ctx,
		opts: options,
	}
}

func WithContext(ctx *cli.Context, opts ...source.Option) source.Source {
	return &cliSource{
		ctx:  ctx,
		opts: source.NewOptions(opts...),
	}
}

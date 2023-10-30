package selector

import (
	"math/rand"
	"time"

	"github.com/wxc/micro/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Random(services []*registry.Service) Next {
	panic(" in Random")
}

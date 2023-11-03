package micro

import "sync"

type service struct {
	once sync.Once
}

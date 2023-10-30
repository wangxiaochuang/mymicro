package source

import (
	"crypto/md5"
	"fmt"
)

func (c *ChangeSet) Sum() string {
	h := md5.New()
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

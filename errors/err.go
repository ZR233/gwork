package errors

import (
	"fmt"
	"github.com/ZR233/goutils/stackx"
)

func HandleRecover(p interface{}) (err error) {
	if p == nil {
		return
	}
	stack := stackx.Stack(1)
	err = fmt.Errorf("%s\n%s", p, stack)
	return
}

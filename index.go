package gwork

import "context"

func New() *WorkPool {
	w := &WorkPool{}
	w.works = map[string]Work{}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	return w
}

package gwork

import "context"

type WorkOptions struct {
	RunImmediately bool //程序启动后立即执行一次
	Description    string
	Name           string
	ReportToCenter bool //执行结果发送至管理中心
	OnError        OnError
}

type WorkPoolOptions struct {
	Prefix     string //任务名称前缀
	ConsulAddr string //若不为空，则向管理中心发送统计信息
}

func NewPool(options *WorkPoolOptions) *WorkPool {
	w := &WorkPool{}
	if options == nil {
		w.WorkPoolOptions = *NewWorkPoolOptions()
	}

	w.works = map[string]Work{}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	return w
}
func NewWorkPoolOptions() *WorkPoolOptions {
	return &WorkPoolOptions{}
}
func NewWorkOptions() *WorkOptions {
	return &WorkOptions{
		RunImmediately: true,
		Description:    "",
	}
}
func NewSchedule() *Schedule {
	return &Schedule{}
}

package job

import (
	"context"
)

// Handler 任务处理器接口
type Handler interface {
	// JobName 返回处理的任务名称
	JobName() string

	// Handle 处理任务
	Handle(ctx context.Context, payload []byte) error
}

// HandlerProvider 是 ServiceProvider 的可选能力接口。
// 实现它的 provider 声明自己拥有的任务处理器；worker 命令启动时跨所有
// provider 收集这些处理器并注册到队列。多个 provider 可各自贡献处理器。
type HandlerProvider interface {
	Jobs() []Handler
}

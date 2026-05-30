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

// A provider declares its job handlers by pushing them with support.ServiceProvider.Add;
// the worker command claims every Handler from the application's contributions.

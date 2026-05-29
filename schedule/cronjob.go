package schedule

import "context"

// CronJob 定义一个自注册的定时任务。
// 实现者通过 Schedule 方法配置调度频率和选项，通过 Run 方法执行业务逻辑。
type CronJob interface {
	// Name 返回任务的唯一标识（如 "workspace:clean-suspended"）。
	Name() string

	// Schedule 配置调度频率和选项。
	// 调度器预先创建 Event 并传入，实现者链式调用配置方法后返回。
	Schedule(event *Event) *Event

	// Run 执行任务逻辑。
	Run(ctx context.Context) error
}

// CronProvider 是 ServiceProvider 的可选能力接口。
// 实现它的 provider 声明自己拥有的定时任务；scheduler 命令启动时跨所有
// provider 收集这些任务并注册到调度器。多个 provider 可各自贡献任务。
type CronProvider interface {
	CronJobs() []CronJob
}

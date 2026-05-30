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

// A provider declares its cron jobs by pushing them with support.ServiceProvider.Add;
// the scheduler command claims every CronJob from the application's contributions.

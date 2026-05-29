package bus

type Listener interface {
	Handlers() map[string]Handler
}

// ListenerProvider 是 ServiceProvider 的可选能力接口。
// 实现它的 provider 声明自己拥有的事件监听器；eventbus 命令启动时跨所有
// provider 收集这些监听器并订阅。多个 provider 可各自贡献监听器。
type ListenerProvider interface {
	Listeners() []Listener
}

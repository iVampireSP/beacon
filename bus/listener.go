package bus

// Listener maps event-name patterns to their handlers. A provider declares its
// listeners by pushing them with support.ServiceProvider.Add; the eventbus
// command claims every Listener from the application's contributions.
type Listener interface {
	Handlers() map[string]Handler
}

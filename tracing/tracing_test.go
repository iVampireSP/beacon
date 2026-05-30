package tracing

import "testing"

// TestGetServiceNoopWhenDisabled verifies tracing is off by default: once the
// provider initializes it from a disabled config, GetService is a no-op (no
// error, nil provider) so commands run without tracing.
func TestGetServiceNoopWhenDisabled(t *testing.T) {
	// Mirrors what the ServiceProvider's Boot does with tracing.enabled=false.
	if _, err := NewTracing(Config{Enabled: false}); err != ErrTracingDisabled {
		t.Fatalf("NewTracing(disabled) = %v, want ErrTracingDisabled", err)
	}

	tp, err := GetService("test-service")
	if err != nil {
		t.Fatalf("GetService with tracing disabled returned error: %v", err)
	}
	if tp != nil {
		t.Fatalf("GetService with tracing disabled returned non-nil provider: %v", tp)
	}
}

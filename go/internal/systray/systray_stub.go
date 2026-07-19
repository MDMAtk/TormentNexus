//go:build !windows

package systray

import (
	"os"
	"github.com/MDMAtk/TormentNexus/internal/eventbus"
)

// Start is a no-op on non-Windows platforms
func Start(eb *eventbus.EventBus) {
	// Headless mode: no system tray UI
}

// NotifyActivity is a no-op on non-Windows platforms
func NotifyActivity(dir string) {
	// Headless mode: do nothing
}

// TriggerFullShutdown is a stub on non-Windows platforms that simply exits the process
func TriggerFullShutdown() {
	// Headless mode: just exit
	os.Exit(0)
}

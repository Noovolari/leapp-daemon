package timer

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	var counter = 0
	scheduledFunction := func() {
		counter++
	}
	timer := NewTimer(1, scheduledFunction)

	var lastCounter = 0
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		if counter <= lastCounter {
			t.Errorf("Timer must call the scheduled function when active")
		}
		lastCounter = counter
	}

	timer.Close()
	lastCounter = counter
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		if counter != lastCounter {
			t.Errorf("Timer must not call the scheduled function when stopped")
		}
	}
}

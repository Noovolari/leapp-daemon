package timer

type TimerCollection struct {
	timers []*timer
}

func NewTimerCollection() *TimerCollection {
	return &TimerCollection{
		timers: make([]*timer, 0),
	}
}

func (timers *TimerCollection) AddTimer(intervalInSeconds int, scheduledFunction ScheduledFunction) {
	timers.timers = append(timers.timers, NewTimer(intervalInSeconds, scheduledFunction))
}

func (timers *TimerCollection) Close() {
	for _, timer := range timers.timers {
		timer.Close()
	}
}

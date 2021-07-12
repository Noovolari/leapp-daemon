package providers

import (
	"leapp_daemon/infrastructure/timer"
	"sync"
)

var timerCollectionSingleton *timer.TimerCollection
var timerCollectionMutex sync.Mutex

func (prov *Providers) GetTimerCollection() *timer.TimerCollection {
	timerCollectionMutex.Lock()
	defer timerCollectionMutex.Unlock()

	if timerCollectionSingleton == nil {
		timerCollectionSingleton = timer.NewTimerCollection()
	}
	return timerCollectionSingleton
}

package metrics

import "sync"

func Reset() {
	isInit = false
	server = nil

	if registry != nil {
		registry.Unregister(concurrentCalls)
		registry.Unregister(totalCalls)
		registry.Unregister(callDuration)
	}

	pubOnce = &sync.Once{}
}

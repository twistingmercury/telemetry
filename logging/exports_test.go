package logging

func SetExitFunction(f func(int)) {
	exitFunc = f
}

var TraceInfo = traceInfo

package logging

func SetExitFunc(f func(int)) {
	exitFunc = f
}

var TraceInfo = traceInfo

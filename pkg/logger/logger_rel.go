//go:build !debug
// +build !debug

package logger

func (l *Logger) getLineFunc(logInfo *LogInfo) {
	logInfo.TrackLine = false
}

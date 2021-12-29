package notifybot

type Logger interface {
	Debug(format string, values ...interface{})
	Info(format string, values ...interface{})
	Warning(format string, values ...interface{})
	Error(format string, values ...interface{})
}

package epilog

var (
	stdio *Logger = New()
)

// Use is an exported function.
func Use() {
	stdio.Use()
}

// End is an exported function.
func End() {
	stdio.End()
}

// ChangeLevel is an exported function.
func ChangeLevel(level logLevel) {
	stdio.ChangeLevel(level)
}

// RegisterAppender is an exported function.
func RegisterAppender(name string, appender Appender) (err error) {
	return stdio.RegisterAppender(name, appender)
}

// Print is an exported function.
func Print(args ...interface{}) {
	stdio.Print(args...)
}

// Debugln is an exported function.
func Debugln(args ...interface{}) {
	stdio.Debugln(args...)
}

// Infoln is an exported function.
func Infoln(args ...interface{}) {
	stdio.Infoln(args...)
}

// Warningln is an exported function.
func Warningln(args ...interface{}) {
	stdio.Warningln(args...)
}

// Fatalln is an exported function.
func Fatalln(args ...interface{}) {
	stdio.Fatalln(args...)
}

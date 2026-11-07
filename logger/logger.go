package logger

type Logger struct {
	level Level
	e     *engine
}

// uses the [DefaultConfig], if modifying anything recommended to copy
// the default config and modify values there and pass that
//
// Do not pass this as a pointer,
// this already holds a pointer to the underlying logging engine,
// so just use this as is.
func New(l Level, opts ...Option) Logger {
	return NewWithConfig(l, &DefaultConfig, opts...)
}

func NewWithConfig(l Level, c *Config, opts ...Option) Logger {
	for _, op := range opts {
		op(c)
	}
	return Logger{
		level: l,
		e:     newEngine(c),
	}
}

func (l Logger) log(lv level, m string, kv []string) {
	if l.e == nil {
		panic("logger: engine is nil, this should not happen, you might be using a empty logger instance without a engine.")
	}
	l.e.log(lv, m, kv)
}

func (l Logger) Trace(m string, kv ...string) {
	l.log(levelTrace, m, kv)
}

func (l Logger) Debug(m string, kv ...string) {
	if l.level <= 1 {
		l.log(levelDebug, m, kv)
	}
}
func (l Logger) Info(m string, kv ...string) {
	if l.level <= 2 {
		l.log(levelInfo, m, kv)
	}
}
func (l Logger) Warn(m string, kv ...string) {
	if l.level <= 3 {
		l.log(levelWarn, m, kv)
	}
}
func (l Logger) Error(m string, kv ...string) {
	if l.level <= 4 {
		l.log(levelError, m, kv)
	}
}

func (l Logger) Fatal(m string, kv ...string) {
	if l.level <= 5 {
		l.log(levelFatal, m, kv)
		if l.e.config.Others.PanicAtFatal {
			panic("fatal log")
		}
	}
}

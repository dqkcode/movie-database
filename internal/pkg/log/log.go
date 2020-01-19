package log

type (
	Logger interface {
		//Infof print info with format.
		Infof(f string, v ...interface{})

		//Debugf print Debug with format.
		Debugf(f string, v ...interface{})

		//Warnf print Warn with format.
		Warnf(f string, v ...interface{})

		//Errorf print Error with format.
		Errorf(f string, v ...interface{})

		//Panicf print Panic with format.
		Panicf(f string, v ...interface{})

		//Info print info.
		Info(v ...interface{})

		//Infof print Debug.
		Debug(v ...interface{})

		//Warn print Warn.
		Warn(v ...interface{})

		//Error print Error.
		Error(v ...interface{})

		//Panic print Panic.
		Panic(v ...interface{})

		//WithField print with field.
		WithField(k, v interface{}) Logger

		//WithFields print with fields.
		WithFields(fields Fields) Logger
	}
	//Fields is alias of map
	Fields = map[string]interface{}

	contextKey string
)

const (
	loggerKey  contextKey = contextKey("logger_key")
	filePrefix            = "file://"
)

var (
	root Logger
)

func Root() Logger {
	if root == nil {
		root = newGlog()
	}
	return root
}

// WithFields return a new logger entry with fields
func WithFields(fields Fields) Logger {
	return Root().WithFields(fields)
}

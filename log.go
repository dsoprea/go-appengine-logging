package log

import (
    "text/template"

    "os"
    "errors"
    "fmt"
    "bytes"

    "golang.org/x/net/context"
)

// Config keys.
const (
    CkLogDefaultFormat = "LogDefaultFormat"
    CkLogDefaultAdapter = "LogDefaultAdapter"
    CkLogDefaultLevel = "LogDefaultLevel"
)

// Config severity integers.
const (
    LevelDebug = iota
    LevelInfo = iota
    LevelWarning = iota
    LevelError = iota
    LevelCritical = iota
)

// Config severity names.
const (
    LevelNameDebug = "debug"
    LevelNameInfo = "info"
    LevelNameWarning = "warning"
    LevelNameError = "error"
    LevelNameCritical = "critical"
)

// Seveirty name->integer map.
var (
    LevelNameMap = map[string]int {
        LevelNameDebug: LevelDebug,
        LevelNameInfo: LevelInfo,
        LevelNameWarning: LevelWarning,
        LevelNameError: LevelError,
        LevelNameCritical: LevelCritical,
    }
)

// Config
var (
    LogDefaultFormat = os.Getenv(CkLogDefaultFormat)
    LogDefaultAdapter = os.Getenv(CkLogDefaultAdapter)
    LogDefaultLevel = os.Getenv(CkLogDefaultLevel)
)

// Errors
var (
    ErrAdapterMakerAlreadyDefined = errors.New("Adapter-maker already defined.")
    ErrLogLevelInvalid = errors.New("Log-level not valid.")
)

var (
    includeFilters = make(map[string]bool)
    useIncludeFilters = true
    excludeFilters = make(map[string]bool)
    useExcludeFilters = true

    makers = make(map[string]*AdapterMaker)
)

// Add global include filter.
func AddIncludeFilter(noun string) {
    includeFilters[noun] = true
    useIncludeFilters = true
}

// Add global exclude filter.
func AddExcludeFilter(noun string) {
    excludeFilters[noun] = true
    useExcludeFilters = true
}

type AdapterMaker interface {
    New() *LogAdapter
}

func AddAdapterMaker(name string, am *AdapterMaker) {
    if _, found := makers[name]; found == true {
        panic(ErrAdapterMakerAlreadyDefined)
    }

    makers[name] = am
}

type LogAdapter interface {
    Criticalf(lc *LogContext, message *string) error
    Debugf(lc *LogContext, message *string) error
    Errorf(lc *LogContext, message *string) error
    Infof(lc *LogContext, message *string) error
    Warningf(lc *LogContext, message *string) error
}

type MessageContext struct {
    Noun *string
    Message *string
}

type LogContext struct {
    Logger *Logger
    Ctx context.Context
}

type Logger struct {
    la *LogAdapter
    t *template.Template
    systemLevel int

    Noun *string
}

func NewLoggerWithAdapter(noun string, la *LogAdapter) *Logger {
    l := &Logger{
        la: la,
        Noun: &noun,
    }

    // Set the level.

    systemLevelName := LogDefaultLevel
    var systemLevel int
    var found bool

    if systemLevelName == "" {
        systemLevel = LevelInfo
    } else if systemLevel, found = LevelNameMap[systemLevelName]; found == false {
        panic(ErrLogLevelInvalid)
    }

    l.systemLevel = systemLevel

    // Set the form.

    format := LogDefaultFormat
    if format == "" {
        format = "{{.Noun}}: {{.Message}}"
    }

    l.SetFormat(format)

    return l
}

func NewLogger(noun string) *Logger {
    var la LogAdapter

    if LogDefaultAdapter != "" {
        am := makers[LogDefaultAdapter]
        la = *(*am).New()
    } else {
        la = AppEngineLogAdapter{}
    }

    return NewLoggerWithAdapter(noun, &la)
}

func (l *Logger) SetFormat(format string) {
    if t, err := template.New("logItem").Parse(format); err != nil {
        panic(err)
    } else {
        l.t = t
    }
}

func (l *Logger) flattenMessage(format *string, args []interface{}) (string, error) {
    m := fmt.Sprintf(*format, args...)
    
    lc := &MessageContext{
        Noun: l.Noun,
        Message: &m,
    }

    var b bytes.Buffer
    if err := l.t.Execute(&b, *lc); err != nil {
        return "", err
    }

    return b.String(), nil
}

func (l *Logger) allowMessage(noun *string) bool {
    if _, found := includeFilters[*noun]; found == true {
        return true
    }

    // If we didn't hit an include filter and we *had* include filters, filter 
    // it out.
    if useIncludeFilters == true {
        return false
    }

    if _, found := excludeFilters[*noun]; found == true {
        return false
    }

    return true
}

func (l *Logger) makeLogContext(ctx context.Context) *LogContext {
    return &LogContext{
        Ctx: ctx,
        Logger: l,
    }
}

func (l *Logger) Criticalf(ctx context.Context, format string, args ...interface{}) error {
    if(l.allowMessage(l.Noun) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return (*l.la).Criticalf(lc, &s)
    }
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) error {
    if l.systemLevel > LevelDebug {
        return nil
    }

    if(l.allowMessage(l.Noun) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return (*l.la).Debugf(lc, &s)
    }
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) error {
    if l.systemLevel > LevelError {
        return nil
    }

    if(l.allowMessage(l.Noun) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return (*l.la).Errorf(lc, &s)
    }
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) error {
    if l.systemLevel > LevelInfo {
        return nil
    }

    if(l.allowMessage(l.Noun) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return (*l.la).Infof(lc, &s)
    }
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) error {
    if l.systemLevel > LevelWarning {
        return nil
    }

    if(l.allowMessage(l.Noun) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return (*l.la).Warningf(lc, &s)
    }
}

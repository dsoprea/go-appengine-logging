package log

import (
    "text/template"

    "os"
    "errors"
    "fmt"
    "bytes"
    "strings"

    "golang.org/x/net/context"
)

// Config keys.
const (
    ckFormat = "LogFormat"
    ckAdapterName = "LogAdapterName"
    ckLevelName = "LogLevelName"
    ckIncludeNouns = "LogIncludeNouns"
    ckExcludeNouns = "LogExcludeNouns"
    ckExcludeBypassLevelName = "LogExcludeBypassLevelName"
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

// Built-in adapter makers
const (
    adapterMakerAppEngine = "appengine"
)

// Other constants
const (
    defaultAdapter = adapterMakerAppEngine
    defaultFormat = "{{.Noun}}: {{.Message}}"
    defaultLevelName = LevelNameInfo
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
    // Alternative format (defaults to .
    format = os.Getenv(ckFormat)

    // Alternative adapter (defaults to "appengine").
    adapterName = os.Getenv(ckAdapterName)

    // Alternative level at which to display log-items (defaults to 
    // "info").
    levelName = os.Getenv(ckLevelName)

    // Configuration-driven comma-separated list of nouns to include. Defaults
    // to empty.
    includeNouns = os.Getenv(ckIncludeNouns)

    // Configuration-driven comma-separated list of nouns to exclude. Defaults 
    // to empty.
    excludeNouns = os.Getenv(ckExcludeNouns)

    // Level at which to disregard exclusion (if the severity of a message 
    // meets or exceed this, always display). Defaults to empty.
    excludeBypassLevelName = os.Getenv(ckExcludeBypassLevelName)
)

// Errors
var (
    ErrAdapterMakerAlreadyDefined = errors.New("adapter-maker already defined")
    ErrLogLevelInvalid = errors.New("log-level not valid")
    ErrFormatEmpty = errors.New("format is empty")
    ErrExcludeLevelNameInvalid = errors.New("exclude bypass-level is invalid")
)

// Other
var (
    includeFilters = make(map[string]bool)
    useIncludeFilters = false
    excludeFilters = make(map[string]bool)
    useExcludeFilters = false

    makers = make(map[string]AdapterMaker)

// TODO(dustin): !! Finish implementing this.
    excludeBypassLevel = -1
)

// Add global include filter.
func AddIncludeFilter(noun string) {
    includeFilters[noun] = true
    useIncludeFilters = true
}

// Remove global include filter.
func RemoveIncludeFilter(noun string) {
    delete(includeFilters, noun)
    if len(includeFilters) == 0 {
        useIncludeFilters = false
    }
}

// Add global exclude filter.
func AddExcludeFilter(noun string) {
    excludeFilters[noun] = true
    useExcludeFilters = true
}

// Remove global exclude filter.
func RemoveExcludeFilter(noun string) {
    delete(excludeFilters, noun)
    if len(excludeFilters) == 0 {
        useExcludeFilters = false
    }
}

type AdapterMaker interface {
    New() LogAdapter
}

type AppEngineAdapterMaker struct {

}

func (aeam AppEngineAdapterMaker) New() LogAdapter {
    return AppEngineLogAdapter{}
}

func AddAdapterMaker(name string, am AdapterMaker) {
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

// TODO(dustin): !! Also populate whether we've bypassed an exception so that 
//                  we can add a template macro to prefix an exclamation of 
//                  some sort.
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

    var systemLevel int
    var found bool
    if systemLevel, found = LevelNameMap[levelName]; found == false {
        panic(ErrLogLevelInvalid)
    }

    l.systemLevel = systemLevel

    // Set the form.

    format := format
    if format == "" {
        panic(ErrFormatEmpty)
    }

    l.SetFormat(format)

    return l
}

func NewLogger(noun string) *Logger {
    var la LogAdapter

    am := makers[adapterName]
    la = am.New()

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

func (l *Logger) allowMessage(noun *string, level int) bool {
    // Preempt the normal filter checks if we can unconditionally allow at a 
    // certain level and we've hit that level.
    if level >= excludeBypassLevel && excludeBypassLevel != -1 {
        return true
    }

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

type LogMethod func(lc *LogContext, message *string) error

func (l *Logger) log(ctx context.Context, level int, lm LogMethod, format string, args []interface{}) error {
    if l.systemLevel > level {
        return nil
    }

    if(l.allowMessage(l.Noun, level) == false) {
        return nil
    }

    if s, err := l.flattenMessage(&format, args); err != nil {
        return err
    } else {
        lc := l.makeLogContext(ctx)
        return lm(lc, &s)
    }
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) error {
    return l.log(ctx, LevelDebug, (*l.la).Debugf, format, args)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) error {
    return l.log(ctx, LevelInfo, (*l.la).Infof, format, args)
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) error {
    return l.log(ctx, LevelWarning, (*l.la).Warningf, format, args)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) error {
    return l.log(ctx, LevelError, (*l.la).Errorf, format, args)
}

func (l *Logger) Criticalf(ctx context.Context, format string, args ...interface{}) error {
    return l.log(ctx, LevelCritical, (*l.la).Criticalf, format, args)
}

func init() {
    if format == "" {
        format = defaultFormat
    }

    if adapterName == "" {
        adapterName = defaultAdapter
    }

    if levelName == "" {
        levelName = defaultLevelName
    }

    AddAdapterMaker(adapterMakerAppEngine, AppEngineAdapterMaker{})

    if includeNouns != "" {
        for _, noun := range strings.Split(includeNouns, ",") {
            AddIncludeFilter(noun)
        }
    }

    if excludeNouns != "" {
        for _, noun := range strings.Split(excludeNouns, ",") {
            AddExcludeFilter(noun)
        }
    }

    if excludeBypassLevelName != "" {
        var found bool
        if excludeBypassLevel, found = LevelNameMap[excludeBypassLevelName]; found == false {
            panic(ErrExcludeLevelNameInvalid)
        }
    }
}

Go's logging is pretty basic by design. Go under AppEngine is even more stripped down. For example, there is no logging struct (`Logger` under traditional Go) and there is no support for prefixes. This package not only equips an AppEngine application with both, but adds include/exclude filters, pluggable logging adapters, and configuration-driven design. This allows you to be able to turn on and off logging for certain packages/files or dependencies as well as to be able to do it from configuration (the YAML files, using `os.Getenv()`).


## Getting Started

The simplest, possible example:

```go
package thispackage

import (
    "golang.org/x/net/context"

    "github.com/dsoprea/go-appengine-logging"
)

var (
    thisfile_log = log.NewLogger("thisfile")
)

func a_cry_for_help(ctx context.Context) {
    thisfile_log.Errorf(ctx, "How big is my problem: %s", "pretty big")
}
```

Notice that we pass in the name of a prefix (what we refer to as a "noun") to `log.NewLogger()`. This is a simple, descriptive name that represents the current body of logic. We recommend that you define a different log for every file at the package level, but it is your choice if you want to go with this methodology, share the same logger over the entire package, define one for each struct, etc..

### Example Output

Example output from a real application (not from the above):

```
2016/09/09 12:57:44 DEBUG: user: User revisiting: [test@example.com]
2016/09/09 12:57:44 DEBUG: context: Session already inited: [DCRBDGRY6RMWANCSJXVLD7GULDH4NZEB6SBAQ3KSFIGA2LP45IIQ]
2016/09/09 12:57:44 DEBUG: session_data: Session save not necessary: [DCRBDGRY6RMWANCSJXVLD7GULDH4NZEB6SBAQ3KSFIGA2LP45IIQ]
2016/09/09 12:57:44 DEBUG: context: Got session: [DCRBDGRY6RMWANCSJXVLD7GULDH4NZEB6SBAQ3KSFIGA2LP45IIQ]
2016/09/09 12:57:44 DEBUG: session_data: Found user in session.
2016/09/09 12:57:44 DEBUG: cache: Cache miss: [geo.geocode.reverse:dhxp15x]
```


## Adapters

This project provides one built-in logging adapter: "appengine", which invokes the default AppEngine logger. If you would like to implement your own logger, just create a struct type that satisfies the LogAdapter interface.

```go
type LogAdapter interface {
    Criticalf(lc *LogContext, message *string) error
    Debugf(lc *LogContext, message *string) error
    Errorf(lc *LogContext, message *string) error
    Infof(lc *LogContext, message *string) error
    Warningf(lc *LogContext, message *string) error
}
```

The *LogContext* struct passed in provides additional information that you may need in order to do what you need to do:

```go
type LogContext struct {
    Logger *Logger
    Ctx context.Context
}
```

Note that *Logger* represents your Logger instance. It exports `Noun *string` in the event you want to discriminate where your log entries go.

Adapter example:

```go
type DummyLogAdapter struct {

}

func (dla DummyLogAdapter) Criticalf(lc *LogContext, message *string) error {
    
}

func (dla DummyLogAdapter) Debugf(lc *LogContext, message *string) error {
    
}

func (dla DummyLogAdapter) Errorf(lc *LogContext, message *string) error {
    
}

func (dla DummyLogAdapter) Infof(lc *LogContext, message *string) error {
    
}

func (dla DummyLogAdapter) Warningf(lc *LogContext, message *string) error {
    
}
```

There are a couple of ways to tell Logger to use a specific adapter:

1. Instead of calling `log.NewLogger(noun string)`, call `log.NewLoggerWithAdapter(noun string, la *LogAdapter)` and provide a struct of your adapter type.
2. Register a factory type for your adapter and set the name of the adapter into your YAML configuration (under `env_variables`).


The factory must satisfy the *AdapterMaker* interface:

```go
type AdapterMaker interface {
    New() LogAdapter
}
```

An example factory and registration of the factory:

```go
type DummyLogAdapterMaker struct {
    
}

func (dlam DummyLogAdapterMaker) New() log.LogAdapter {
    return DummyLogAdapter{}
}
```

We then recommending registering it from the `init()` function of the fiel that defines the maker type:

```go
func init() {
    log.AddAdapterMaker("dummy", DummyLogAdapterMaker{})
}
```

We discuss how to then reference the adapter-maker from configuration in the "Configuration" section below.


### Filters

We support the ability to exclusively log for a specific set of nouns (we'll exclude any not specified):

```go
log.AddIncludeFilter("nountoshow1")
log.AddIncludeFilter("nountoshow2")
```

Depending on your needs, you might just want to exclude a couple and include the rest:

```go
log.AddExcludeFilter("nountohide1")
log.AddExcludeFilter("nountohide2")
```

We'll first hit the include-filters. If it's in there, we'll forward the log item to the adapter. If not, and there is at least one include filter in the list, we won't do anything. If the list of include filters is empty but the noun appears in the exclude list, we won't do anything.


#### Footnote

It is a good convention to exclude the nouns of any library you are writing whose logging you do not want to generally be aware of unless you are debugging. You might call `AddExcludeFilter()` from the `init()` function at the bottom of those files unless there is some configuration variable, such as "(LibraryNameHere)DoShowLogging", that has been defined and set to TRUE.


### Configuration

The following configuration keys can be used to pre-configure your logging directly from your configuration (the AppEngine YAML files):

- *LogFormat*: The default format used to build the message that gets sent to the adapter. It is assumed that the adapter already prefixes the message with time and log-level (since the default AppEngine logger does). The default value is: `{{.Noun}}:{{if eq .ExcludeBypass true}} [BYPASS]{{end}} {{.Message}}`. The available tokens are "Noun", "ExcludeBypass", and "Message".
- *LogAdapterName*: The name of the adapter to use when NewLogger() is called.
- *LogLevelName*: The priority-level of messages permitted to be logged (all others will be discarded). By default, it is "info". Other levels are: "debug", "warning", "error", "critical"
- *LogIncludeNouns*: Comma-separated list of nouns to log for. All others will be ignored.
- *LogExcludeNouns*: Comma-separated list on nouns to exclude from logging.
- *LogExcludeBypassLevelName*: The log-level at which we will show logging for nouns that have been excluded. Allows you to hide excessive, unimportant logging for nouns but to still see their warnings, errors, etc...

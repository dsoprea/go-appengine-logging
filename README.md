[![Build Status](https://travis-ci.org/dsoprea/go-appengine-logging.svg?branch=master)](https://travis-ci.org/dsoprea/go-appengine-logging)

## Overview

This project uses the [go-logging](https://github.com/dsoprea/go-logging) project to provide enhanced logging under AppEngine. This includes the availability of stacktraces.

## Example

Usage:

```go
package app

import (
    "golang.org/x/net/context"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-appengine-logging"
)

var (
    // You should create one of these at the top of every file and name them so 
    // it's clear which file they represent.
    thisfileLogger = log.NewLogger("app.thisfile")
)

func someCall(ctx context.Context) {
    thisfileLogger.Debugf(ctx, "Test message.")
    thisfileLogger.Infof(ctx, "Test message.")
    thisfileLogger.Warningf(ctx, "Test message.")
    
    err := e.New("some error")
    thisfileLogger.Errorf(ctx, err, "Test message.")
}

// Do this in whichever file initializes your application. Don't do this in 
// libraries.
func init() {
    aam := aelog.NewAppengineAdapterMaker()
    log.AddAdapterMaker("appengine", aam)
}
```

For more information, see the documentation for [go-logging](https://github.com/dsoprea/go-logging).

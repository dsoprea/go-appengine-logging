[![Build Status](https://travis-ci.org/dsoprea/go-appengine-logging.svg?branch=master)](https://travis-ci.org/dsoprea/go-appengine-logging)

## Overview

This project uses the [go-logging](https://github.com/dsoprea/go-logging) project to provide enhanced logging under AppEngine. This includes the availability of stacktraces.


## Example

Usage:

```go
package app

import (
    e "errors"

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

Example application output:

```
2016/10/01 18:42:25 DEBUG: common.geographic: Geocode result (0) address component (5): [US] [United States] [[country political]]
2016/10/01 18:42:25 DEBUG: common.geographic: Geocode result (0) address component (6): [33436] [33436] [[postal_code]]
2016/10/01 18:42:25 DEBUG: common.geographic: Geocode result (0) address component (7): [8616] [8616] [[postal_code_suffix]]
2016/10/01 18:42:25 DEBUG: app: Flushing session: [SessionData<JW6NOY236NSOJJRXZDG73MRGSJHLHR34JCEHGONMM46N66LIDBOQ>]
INFO     2016-10-01 18:42:25,287 module.py:788] default: "GET /api/1/geo/place/geocode/reverse?lat=26.562653899999997&lon=-80.1022059 HTTP/1.1" 200 103
2016/10/01 18:42:25 DEBUG: data.user: Retrieving user-account: [test@example.com]
2016/10/01 18:42:25 DEBUG: data.user: Found and returning.
2016/10/01 18:42:25 DEBUG: user: User login: [test@example.com]
2016/10/01 18:42:27 DEBUG: common.geographic: Caching place: [ChIJb5J6W0Mn2YgR5MuqicTF5PA] [D.M.T. Preservations LLC] [geo.places.nearby.entity:ChIJb5J6W0Mn2YgR5MuqicTF5PA]
2016/10/01 18:42:27 DEBUG: common.geographic: Caching place: [ChIJq87Z0GAn2YgR9u32K0j_ur4] [Napoli Ristorante Pizzeria] [geo.places.nearby.entity:ChIJq87Z0GAn2YgR9u32K0j_ur4]
2016/10/01 18:42:27 DEBUG: common.geographic: Caching place: [ChIJu_N2xWAn2YgR62tDpH3zjqo] [Sabai Thai Restaurant] [geo.places.nearby.entity:ChIJu_N2xWAn2YgR62tDpH3zjqo]
```

For more information on usage, please see the documentation for [go-logging](https://github.com/dsoprea/go-logging).


## Backwards Compatibility

A split occurred early in the life of [go-appengine-logging](https://github.com/dsoprea/go-appengine-logging). Most of the functionality was moved to [go-logging](https://github.com/dsoprea/go-logging). Though this was the original project, and it had a limited amount of functionality in the beginning, once it got implemented and experienced some natural growth it became very, very useful. However, it, simultaneously, became non-trivial enough that we weren't thrilled about duplicating it in order to have the same benefits in a non-AppEngine-specific project.

Without getting into too much unnecessary detail, the following decisions were made:

1. The main project would be *go-logging*. **Any existing references to *go-appengine-logging* in any projects that depend on this would have to be updated.**
2. Because the *go-logging* code could no longer be intrinsically aware of *go-appengine-logging* (because an AE project has different requirements and expects a different environment), **applications that use *go-logging* must specifically import *go-appengine-logging* and register the one with the other**.

This will break anyone that is not vendoring *go-appengine-logging*, but, as that was an AppEngine-specific project and AE projects are usually predisposed to vendoring everything, this won't be an issue unless you update. Plus, that project is still pretty recent and adoption is still just ramping up.

Sorry for any inconvenience. The original project completely changed the logging and debugging experience for AppEngine and I wanted to bring that over to general Go development.

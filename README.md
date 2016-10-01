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

For more information, see the documentation for [go-logging](https://github.com/dsoprea/go-logging).

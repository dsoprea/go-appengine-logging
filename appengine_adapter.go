package log

import (
    aelog "google.golang.org/appengine/log"
)

type AppEngineLogAdapter struct {

}

// TODO(dustin): !! Fix these to use pointer receivers.

func (ael AppEngineLogAdapter) Criticalf(lc *LogContext, message *string) error {
    aelog.Criticalf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Debugf(lc *LogContext, message *string) error {
    aelog.Debugf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Errorf(lc *LogContext, message *string) error {
    aelog.Errorf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Infof(lc *LogContext, message *string) error {
    aelog.Infof(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Warningf(lc *LogContext, message *string) error {
    aelog.Warningf(lc.Ctx, *message)

    return nil
}

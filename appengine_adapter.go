package log

import (
    "google.golang.org/appengine/log"
)

type AppEngineLogAdapter struct {

}

// TODO(dustin): !! Fix these to use pointer receivers.

func (ael AppEngineLogAdapter) Criticalf(lc *LogContext, message *string) error {
    log.Criticalf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Debugf(lc *LogContext, message *string) error {
    log.Debugf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Errorf(lc *LogContext, message *string) error {
    log.Errorf(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Infof(lc *LogContext, message *string) error {
    log.Infof(lc.Ctx, *message)

    return nil
}

func (ael AppEngineLogAdapter) Warningf(lc *LogContext, message *string) error {
    log.Warningf(lc.Ctx, *message)

    return nil
}

package aelog

import (
    a "google.golang.org/appengine/log"

    "github.com/dsoprea/go-logging"
)

type AppengineLogAdapter struct {

}

func (ala *AppengineLogAdapter) Debugf(lc *log.LogContext, message *string) error {
    a.Debugf(lc.Ctx, *message)

    return nil
}

func (ala *AppengineLogAdapter) Infof(lc *log.LogContext, message *string) error {
    a.Infof(lc.Ctx, *message)

    return nil
}

func (ala *AppengineLogAdapter) Warningf(lc *log.LogContext, message *string) error {
    a.Warningf(lc.Ctx, *message)

    return nil
}

func (ala *AppengineLogAdapter) Errorf(lc *log.LogContext, message *string) error {
    a.Errorf(lc.Ctx, *message)

    return nil
}


type AppengineAdapterMaker struct {

}

func NewAppengineAdapterMaker() *AppengineAdapterMaker {
    return new(AppengineAdapterMaker)
}

func (aam AppengineAdapterMaker) New() log.LogAdapter {
    return new(AppengineLogAdapter)
}

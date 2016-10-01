package aelog

import (
    "testing"

    "github.com/dsoprea/go-logging"

    "golang.org/x/net/context"
)

func TestAppengine(t *testing.T) {
    ecp := log.NewEnvironmentConfigurationProvider()
    log.LoadConfiguration(ecp)

    log.ClearAdapters()

    aam := NewAppengineAdapterMaker()
    log.AddAdapterMaker("appengine", aam)

    an := log.GetDefaultAdapterName()
    if an != "appengine" {
        t.Error("Adapter was not properly registered.")
    }

    // Usually this is either useless or fundamentally broken in the context of
    // AppEngine, but it works here and we won't otherwise have a useful 
    // context since we're not running in response to a request.
    ctx := context.Background()

    l := log.NewLoggerWithAdapter("appengine_test", "appengine")
    l.Debugf(ctx, "Test message.")
    l.Infof(ctx, "Test message.")
    l.Warningf(ctx, "Test message.")
    l.Errorf(ctx, "Test message.")
}

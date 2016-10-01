package aelog

import (
    "testing"
    e "errors"

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

    // We can't actually test the logging calls unless we establish a full GAE 
    // environment. We'll leave this as an exercise for later:
    //
    // http://stackoverflow.com/questions/24614599/compile-app-engine-application-in-travis
    // http://orcaman.blogspot.com/2014/09/ci-when-githubtravis-meet-gogae.html
    // https://github.com/golang/appengine/pull/5/files
}

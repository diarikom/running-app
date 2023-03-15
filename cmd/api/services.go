package main

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/pkg/errors"
	"reflect"
)

func initServices(app *api.Api) {
	// Get logger
	logger := app.Logger

	// Reflect services
	rv := reflect.ValueOf(app.Services).Elem()

	// Iterate fields
	for i := 0; i < rv.NumField(); i++ {
		// Get field interface
		f := rv.Field(i).Interface()
		svcName := rv.Type().Field(i).Name

		// Cast interface to
		svc, ok := f.(api.ServiceInitiator)
		if !ok {
			panic(fmt.Errorf("api: Services.%s does not implement ServiceInitiator interface", svcName))
		}

		// Init service
		err := svc.Init(app)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("api: Services.%s error while init", svcName)))
		}

		logger.Debugf("Services.%s initiated", svcName)
	}
}

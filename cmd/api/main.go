package main

import (
	"flag"
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nlogrus"
	"github.com/gorilla/mux"
	"net/http"
)

var log nlog.Logger

func init() {
	// Create Logger
	logger := nlogrus.NewConsoleLogger(nlogrus.GetLoggerOptEnv())

	// Register logger
	nlog.Register(logger)

	// Get logger instance
	log = nlog.Get()
}

func main() {
	// Init commands
	appFlags := initAppFlags()

	// Parse commands
	flag.Parse()

	// If help command is called, show help
	if *appFlags.CmdShowHelp {
		showHelp()
		return
	}

	// If version command is called, show version
	if *appFlags.CmdShowVersion {
		showVersion()
		return
	}

	app := boot(appFlags)
	start(&app)
}

func start(app *api.Api) {
	// Init handlers
	handlers := initHandlers(app)

	// Init api router
	router := mux.NewRouter()
	apiRouter := nhttp.NewApiRouter(nhttp.RouterOpt{
		RootRouter: router,
		BasePath:   app.BaseUrl.Path,
		Logger:     app.Logger,
	})

	// Register middlewares
	initMiddlewares(apiRouter, app.Services)

	// Init handlers
	initRoutes(apiRouter, handlers)

	log.Debugf("Boot Time: %s", app.Uptime())
	log.Infof("%s API is listening to port %s", AppName, app.Port)
	log.Infof("%s Started. Base URL: %s", AppName, app.BaseUrl.String())
	err := http.ListenAndServe(":"+app.Port, router)
	if err != nil {
		panic(fmt.Errorf("%s: an error occurred while listening to requests (%s", AppSlug, err))
	}
}

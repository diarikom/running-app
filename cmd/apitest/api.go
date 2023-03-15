/// Package that implements API contract for testing purpose

package apitest

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/mocks"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ncore"
	"github.com/diarikom/running-app/running-app-api/pkg/nlogrus"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

type Api struct {
	*api.Api
}

func InitApi() Api {
	// Init logger
	logger := nlogrus.NewConsoleLogger(nlogrus.GetLoggerOptEnv())

	// Init core
	bootOpt := loadBootOptions()
	core := ncore.Boot(bootOpt)

	// Init database
	db := initDb(core.Config)

	// Init error util components
	errUtil := api.NewErrorUtil(bootOpt.ErrorCodesPath, true)
	logger.Debug("Components.Errors initiated")

	// Init id generator components
	idGen := api.NewSnowflakeGenerator(bootOpt.NodeNo)
	logger.Debug("Components.Id initiated")

	// Init app
	app := api.Api{
		BaseUrl: &url.URL{},
		Datasources: &api.Datasources{
			Db: db,
		},
		Components: &api.Components{
			Errors:    errUtil,
			Id:        idGen,
			JWTIssuer: &mocks.JWTIssuerComponent{},
			Mailer:    &mocks.MailerComponent{},
			Facebook:  &mocks.FacebookProviderComponent{},
		},
		Logger: logger,
		Core:   core,
	}

	return Api{Api: &app}
}

func (a *Api) MustInitService(name string, iSvc interface{}) {
	svc, ok := iSvc.(api.ServiceInitiator)
	if !ok {
		panic(fmt.Errorf("apitest: Service %s does not implement ServiceInitiator", name))
	}

	err := svc.Init(a.Api)
	if err != nil {
		panic(errors.Wrap(err, "apitest: error occurred while init Service "+name))
	}
}

func (a *Api) IgnoreDbExec(query string, args ...interface{}) {
	db := a.Datasources.Db.Conn

	result, err := db.Exec(query, args...)
	if err != nil {
		a.Logger.Warnf("failed to execute query. Query = %s, Error = %s", query, err)
	} else {
		a.Logger.Debugf("Result = %+v", result)
	}
}

func loadBootOptions() ncore.BootOpt {
	// Determine working dir
	workDir, ok := os.LookupEnv("WORKING_DIR")
	if !ok {
		workDir = "."
	}

	// Determine config path
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = workDir + "/config.yml"
	}

	// Determine error codes path
	errCodesPath, ok := os.LookupEnv("ERROR_CODES_PATH")
	if !ok {
		errCodesPath = workDir + "/error-codes.yml"
	}

	// Load boot options
	bootOpt := ncore.BootOpt{
		Environment:    ncore.TestingEnvironment,
		Dir:            workDir,
		ConfigPath:     configPath,
		ConfigRequired: []string{},
		ErrorCodesPath: errCodesPath,
		NodeNo:         0,
	}

	return bootOpt
}

func initDb(config *viper.Viper) *nsql.SqlDatabase {
	// Load db config
	var dbConf nsql.Config
	err := config.UnmarshalKey("datasources.db", &dbConf)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to retrieve config for Datasources.Db (%s)", err))
	}

	// Init db
	db, err := nsql.NewSqlDatabase(dbConf)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to init Datasources.Db (%s)", err))
	}

	return db
}

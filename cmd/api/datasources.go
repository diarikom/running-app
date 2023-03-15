package main

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/ns3"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/spf13/viper"
)

func initDatasources(config *viper.Viper) api.Datasources {
	db := initDb(config)
	log.Debug("Datasources.Db initiated")

	asset := initS3(config)
	log.Debug("Datasources.Asset initiated")

	return api.Datasources{
		Db:    db,
		Asset: asset,
	}
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

func initS3(config *viper.Viper) api.S3Provider {
	// Load config
	opt := ns3.MinioOpt{
		Endpoint:        config.GetString(api.ConfAssetEndpoint),
		AccessKeyId:     config.GetString(api.ConfAssetAccessKeyId),
		SecretAccessKey: config.GetString(api.ConfAssetSecretAccessKey),
		UseSSL:          config.GetBool(api.ConfAssetUseSSL),
		BucketName:      config.GetString(api.ConfAssetBucketName),
		Region:          config.GetString(api.ConfAssetRegion),
	}

	// Init S3
	s3, err := ns3.NewMinio(opt)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to init Datasources.Asset (%s)", err))
	}

	return s3
}

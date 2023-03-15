package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type userDiffer struct {
	profile *nsql.Differ
}

func initUserDiffer() userDiffer {
	p := nsql.PrepareDiffer(nsql.DifferOpt{
		Sample:    model.UserProfile{},
		TableName: "user_profile",
	})

	return userDiffer{profile: p}
}

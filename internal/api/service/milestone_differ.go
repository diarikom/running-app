package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type milestoneDiffer struct {
	userChallenge *nsql.Differ
}

func initMilestoneDiffer() milestoneDiffer {
	uc := nsql.PrepareDiffer(nsql.DifferOpt{
		Sample:    model.UserChallenge{},
		TableName: "user_challenge",
	})

	return milestoneDiffer{userChallenge: uc}
}

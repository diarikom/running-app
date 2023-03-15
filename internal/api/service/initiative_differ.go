package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
)

type initiativeDiffer struct {
	donation *nsql.Differ
}

func initInitiativeDiffer() initiativeDiffer {
	donationDiffer := nsql.PrepareDiffer(nsql.DifferOpt{
		Sample:    model.Donation{},
		TableName: "donation",
	})

	return initiativeDiffer{donation: donationDiffer}
}

package api

import (
	"github.com/diarikom/running-app/running-app-api/internal/pkg/nfacebook"
	"github.com/diarikom/running-app/running-app-api/internal/pkg/njwt"
	"github.com/diarikom/running-app/running-app-api/pkg/nmailgun"
)

type JWTIssuerComponent interface {
	New(opt njwt.ClaimOpt) (*njwt.Token, error)
	Verify(input string) (*njwt.Claim, error)
}

type MailerComponent interface {
	GetSender(name string) string
	GetDefaultSender() string
	Send(opt nmailgun.SendOpt) error
}

type FacebookProviderComponent interface {
	GetUrl(path string) string
	InspectToken(token string) (*nfacebook.TokenData, error)
}

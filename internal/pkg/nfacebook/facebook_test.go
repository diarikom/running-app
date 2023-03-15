package nfacebook

import (
	"os"
	"testing"
)

var fbProvider *Provider

func TestMain(m *testing.M) {
	// Get app id and app secret from env
	appId := os.Getenv("TEST_FB_APP_ID")
	appSecret := os.Getenv("TEST_FB_APP_SECRET")

	// Init provider
	var err error
	fbProvider, err = NewProvider(ProviderOpt{
		AppId:     appId,
		AppSecret: appSecret,
	})
	if err != nil {
		panic(err)
	}

	// Run test
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestInspectTokenValidUser(t *testing.T) {
	// Get token from env
	userToken := os.Getenv("TEST_FB_USER_TOKEN")
	if userToken == "" {
		t.Error("Env TEST_FB_USER_TOKEN must not empty")
		return
	}

	// Inspect token
	tokenData, err := fbProvider.InspectToken(userToken)
	if err != nil {
		t.Error(err)
		return
	}

	// Check Type
	if tokenData.Type != EntityUser {
		t.Errorf("token must be a user token. Type: %s", tokenData.Type)
		return
	}

	// Check IsValid
	if !tokenData.IsValid {
		t.Errorf("token must valid.")
		return
	}

	t.Logf("token is a valid user token. UserId: %s", tokenData.UserId)
}

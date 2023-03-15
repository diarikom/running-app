package nfacebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// EntityTypes
	EntityUser = "USER"

	// Api Version
	graphApiVersion = "v6.0"

	// Path
	inspectTokenPath = "/debug_token"
)

type ProviderOpt struct {
	AppId     string
	AppSecret string
}

func NewProvider(opt ProviderOpt) (*Provider, error) {
	// Check app id and app secret
	if opt.AppId == "" || opt.AppSecret == "" {
		return nil, fmt.Errorf("nfacebook: AppId and AppSecret is required")
	}

	// Generate base url
	baseUrl := "https://graph.facebook.com/" + graphApiVersion
	appToken := opt.AppId + "|" + opt.AppSecret

	// Create provider instance
	p := Provider{
		AppId:     opt.AppId,
		AppSecret: opt.AppSecret,
		baseUrl:   baseUrl,
		appToken:  appToken,
	}

	return &p, nil
}

type Provider struct {
	AppId     string
	AppSecret string
	// Private Fields
	baseUrl  string
	appToken string
}

func (p *Provider) GetUrl(path string) string {
	return fmt.Sprintf("%s/%s?access_token=%s", p.baseUrl, path, p.appToken)
}

func (p *Provider) InspectToken(token string) (*TokenData, error) {
	// Generate path url
	u := p.GetUrl(inspectTokenPath)

	// Add input token
	u += "&input_token=" + token

	// Make request to facebook API
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("nfacebook: unable to request InspectToken to Facebook API (%s)", err)
	}
	defer closeResp(resp)

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nfacebook: unable to read InspectToken response (%s)", err)
	}

	// Parse response json
	var respBody struct {
		Data InspectTokenResp `json:"data"`
	}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, fmt.Errorf("nfacebook: unable to unmarshal InspectToken response (%s)", err)
	}

	// If response error, return error
	if errResp := respBody.Data.Error; errResp != nil {
		return nil, fmt.Errorf("nfacebook: received error on InspectToken (%d, %s)", errResp.Code,
			errResp.Message)
	}

	return &respBody.Data.TokenData, nil
}

func closeResp(resp *http.Response) {
	_ = resp.Body.Close()
}

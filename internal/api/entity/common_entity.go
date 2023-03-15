package entity

import "github.com/diarikom/running-app/running-app-api/pkg/nsql"

type ExternalURL struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func (e *ExternalURL) Scan(src interface{}) error {
	return nsql.ScanJSON(src, e)
}

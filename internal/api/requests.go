package api

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nstr"
	"net/url"
)

const (
	KeySkip  = "skip"
	KeyLimit = "limit"
)

func Pagination(q url.Values) (skip int64, limit int8) {
	skip = nstr.ParseInt64(q.Get(KeySkip), SkipDefault)
	limit = nstr.ParseInt8(q.Get(KeyLimit), LimitDefault)
	return
}

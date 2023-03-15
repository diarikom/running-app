package api

import "io"

const (
	AssetsPublicScope = iota
)

type S3Provider interface {
	Upload(file io.Reader, contentType, dest string, scope int) error
}

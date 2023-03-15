package nhttp

import (
	"mime/multipart"
	"path/filepath"
)

// MultipartFile represents structure of multipart file in HTTP request
type MultipartFile struct {
	multipart.File
	Header   *multipart.FileHeader
	MimeType string
}

// Extension returns extension of a filename
func (f *MultipartFile) Extension() string {
	return filepath.Ext(f.Header.Filename)
}

// Rename update file name
func (f *MultipartFile) Rename(newName string) string {
	n := newName + f.Extension()
	f.Header.Filename = n
	return n
}

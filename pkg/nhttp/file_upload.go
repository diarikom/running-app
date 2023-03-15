package nhttp

import (
	"errors"
	"io"
	"net/http"
)

const (
	// Max File Size
	MaxFileSizeImage = 2097152 // 2 MB
	// Default Parameter Keys
	DefaultKeyFile = "file"
)

var (
	// Mime Types
	MimeTypesImage = []string{"image/jpeg", "image/png"}
	// Error Messages
	ErrMimeTypeNotAccepted = errors.New("nhttp: mime type not accepted")
	ErrFileTooLarge        = errors.New("nhttp: request body too large")
)

// UploadRule represent structure for file upload rule and validation
type UploadRule struct {
	Key       string
	MaxSize   int64
	MimeTypes []string
}

// NewImageUploadRules return image with default upload rules
func NewImageUploadRules(keys ...string) (rules []UploadRule) {
	// If keys is unset, add default key file
	if len(keys) == 0 {
		keys = []string{DefaultKeyFile}
	}
	// Create rules
	for _, key := range keys {
		r := UploadRule{key, MaxFileSizeImage, MimeTypesImage}
		rules = append(rules, r)
	}
	// Return rules
	return rules
}

// GetFile retrieve, validate and return file from http.Request
func GetFile(r *http.Request, key string, maxSize int64, mimeTypes []string) (result MultipartFile, err error) {
	// If multipart form has not been parsed, parse with max size limit
	if r.MultipartForm == nil {
		// Parse file by maximum size
		err := r.ParseMultipartForm(maxSize)
		if err != nil {
			return result, err
		}
	}
	// Get file
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return result, err
	}
	defer closeFile(file)
	// If file size too big, return error
	if fileHeader.Size > maxSize {
		return result, ErrFileTooLarge
	}
	// Get mimeType
	mimeType := fileHeader.Header.Get("Content-Type")
	// If mime types is defined, validate mime type
	if mimeTypes != nil && len(mimeTypes) > 0 {
		// Init found flag
		isAccepted := false
		for _, v := range mimeTypes {
			if mimeType == v {
				isAccepted = true
				break
			}
		}
		// If mime type is not mat
		if !isAccepted {
			return result, ErrMimeTypeNotAccepted
		}
	}
	// Return uploaded file
	return MultipartFile{
		File:     file,
		Header:   fileHeader,
		MimeType: mimeType,
	}, nil
}

// GetImages get images in bulk
func GetImages(r *http.Request, config ...UploadRule) (map[string]MultipartFile, error) {
	// If keys is nil or len 0, set default key to file
	if len(config) == 0 {
		config = NewImageUploadRules()
	}
	// Get key length
	l := int64(len(config))
	// If multipart form has not been parsed, parse with max size limit multiplied by no of image file to be uploaded
	if r.MultipartForm == nil {
		// get maximum upload size
		var maxUploadSize int64
		for _, cfg := range config {
			maxUploadSize += cfg.MaxSize
		}
		// Parse file by maximum size
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			return nil, err
		}
	}
	// Init map
	files := make(map[string]MultipartFile, l)
	// Get from keys
	for _, cfg := range config {
		f, err := GetFile(r, cfg.Key, cfg.MaxSize, cfg.MimeTypes)
		if err != nil {
			return nil, err
		}
		// Insert to map
		files[cfg.Key] = f
	}
	// Return
	return files, nil
}

func closeFile(c io.Closer) {
	if err := c.Close(); err != nil {
		// TODO: Print Log
	}
}

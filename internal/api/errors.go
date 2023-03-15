package api

import (
	"fmt"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Errors struct {
	codes map[string]nhttp.Error
	debug bool
}

func (e *Errors) New(code string) nhttp.Error {
	// Get codes
	c, ok := e.codes[code]

	// If not found, return internal server error
	if !ok {
		return nhttp.ErrInternalServer
	}

	// TODO: Trace new error call if in debug mode
	// Copy error
	err := nhttp.Error{
		Status:  c.Status,
		Code:    code,
		Message: c.Message,
	}

	return err
}

func NewErrorUtil(filePath string, debug bool) *Errors {
	errorCodes := loadErrorCodes(filePath)

	return &Errors{
		codes: errorCodes,
		debug: debug,
	}
}

// loadErrorCodes load error codes from file. Error Codes must be YAML formatted
func loadErrorCodes(filePath string) (codes map[string]nhttp.Error) {
	// Read file
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		// If file not exist, panic
		panic(fmt.Errorf("running-app-api: unable to locate error codes file in %s", filePath))
	}
	// Parse error codes file
	err = yaml.Unmarshal(bytes, &codes)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to read error codes file in %s", filePath))
	}
	return codes
}

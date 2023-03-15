package ngrule

import "fmt"

// Error codes constant
const (
	ErrNoNextLogicOp = iota
	ErrActionUnsupported
	ErrConditionUnsupported
	ErrParamUnsupported
	ErrInvalidIntValue
	ErrInvalidFloatValue
	ErrInvalidBooleanValue
	ErrSqlScan
	ErrParamFactFinderNotRegistered
	ErrParamFactFinderExecFail
)

type RuleEngineError struct {
	Code int
}

func (e *RuleEngineError) Error() string {
	return fmt.Sprintf("ErrCode: %d", e.Code)
}

func NewError(code int, msg string) error {
	err := RuleEngineError{
		Code: code,
	}
	return fmt.Errorf("ngrule: %s (%w)", msg, &err)
}

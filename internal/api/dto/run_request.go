package dto

type RunSessionReq struct {
	SessionStarted int64   `json:"session_started"`
	SessionEnded   int64   `json:"session_ended"`
	TimeElapsed    int     `json:"time_elapsed"`
	Distance       int     `json:"distance"`
	Speed          float64 `json:"speed"`
	StepCount      int     `json:"step_count"`
}

type RunStatusSyncReq struct {
	Status int `json:"status"`
}

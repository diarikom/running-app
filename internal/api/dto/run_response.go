package dto

type RunSessionHistoryResp struct {
	RunSessions []RunSessionHistoryItem `json:"run_sessions"`
	Count       int                     `json:"count"`
}

type RunSessionHistoryItem struct {
	Id             string  `json:"id"`
	SessionStarted int64   `json:"session_started"`
	SessionEnded   int64   `json:"session_ended"`
	TimeElapsed    int     `json:"time_elapsed"`
	Distance       int     `json:"distance"`
	Speed          float64 `json:"speed"`
	StepCount      int     `json:"step_count"`
	SyncStatusId   int     `json:"sync_status"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
	Version        int     `json:"version"`
}

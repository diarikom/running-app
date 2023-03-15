package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/dto"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"time"
)

type Run struct {
	IdGen            *api.SnowflakeGen
	Errors           *api.Errors
	Logger           nlog.Logger
	RunRepository    api.RunRepository
	MilestoneService api.MilestoneService
}

func (r *Run) Init(app *api.Api) error {
	// Init run service
	r.IdGen = app.Components.Id
	r.Errors = app.Components.Errors
	r.Logger = app.Logger
	r.RunRepository = NewRunRepository(app.Datasources.Db, app.Logger)
	r.MilestoneService = app.Services.MilestoneService
	return nil
}

func (r Run) GetRunSessionHistory(userId string, skip int64, limit int8) (resp *dto.RunSessionHistoryResp, err error) {
	// Find run session history
	sessions, err := r.RunRepository.FindRunSessionHistory(userId, limit, skip)
	if err != nil {
		r.Logger.Error("unable find run session history", err)
		return nil, err
	}

	// Count total available run session history
	count, err := r.RunRepository.CountRunSessionHistory(userId)
	if err != nil {
		r.Logger.Error("unable count run session", err)
		return nil, err
	}

	// Copy run session history to entity
	items := make([]dto.RunSessionHistoryItem, len(sessions))
	for k, v := range sessions {
		items[k] = dto.RunSessionHistoryItem{
			Id:             v.Id,
			SessionStarted: v.SessionStarted.Unix(),
			SessionEnded:   v.SessionEnded.Unix(),
			TimeElapsed:    v.TimeElapsed,
			Distance:       v.Distance,
			Speed:          v.Speed,
			StepCount:      v.StepCount,
			SyncStatusId:   v.SyncStatusId,
			CreatedAt:      v.CreatedAt.Unix(),
			UpdatedAt:      v.UpdatedAt.Unix(),
			Version:        v.Version,
		}
	}

	// Create response result
	resp = &dto.RunSessionHistoryResp{
		RunSessions: items,
		Count:       count,
	}

	return resp, nil
}

func (r Run) NewRunSession(userId string, req *dto.RunSessionReq) error {
	// Validate required fields
	if req.SessionStarted == 0 ||
		req.SessionEnded == 0 ||
		req.TimeElapsed == 0 ||
		req.Distance == 0 ||
		req.Speed == 0 {

		return nhttp.ErrBadRequest
	}

	// Init timestamp
	timestamp := time.Now()

	// Create run session model
	runSession := model.RunSession{
		Id:             r.IdGen.New(),
		UserId:         userId,
		SessionStarted: time.Unix(req.SessionStarted, 0),
		SessionEnded:   time.Unix(req.SessionEnded, 0),
		TimeElapsed:    req.TimeElapsed,
		Distance:       req.Distance,
		Speed:          req.Speed,
		StepCount:      req.StepCount,
		SyncStatusId:   api.RunSummaryStored,
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
	}

	// Generate Session
	err := r.RunRepository.InsertRunSession(runSession)
	if err != nil {
		r.Logger.Error("unable to persist run session", err)
		return err
	}

	// Trigger check achieved challenge
	err = r.MilestoneService.TriggerCheckChallengeAchieved(dto.UserChallengeReq{
		UserId:    userId,
		Timestamp: time.Now(),
	})
	if err != nil {
		r.Logger.Error("failed to trigger check achieved challenge. UserId = "+userId, err)
	}

	return nil
}

func (r Run) UpdateRunSyncStatus(id, userId string, status int) error {
	// Validate status
	if status != api.RunSummaryStored &&
		status != api.RunDetailsStored {

		return nhttp.ErrBadRequest
	}

	// Persist sync status update
	err := r.RunRepository.UpdateRunSyncStatus(id, userId, status)
	if err != nil {
		r.Logger.Error("unable to persist run sync status update", err)
		return err
	}

	return err
}

package dto

type MilestoneResp struct {
	Milestone *CurrentMilestoneResp `json:"milestone"`
}

type CurrentMilestoneResp struct {
	Id             string                    `json:"id"`
	PeriodStart    int64                     `json:"period_start"`
	PeriodEnd      int64                     `json:"period_end"`
	PeriodTzOffset int                       `json:"period_tz_offset"`
	Status         int                       `json:"status"`
	Challenges     []MilestoneChallengesResp `json:"challenges"`
}

type MilestoneChallengesResp struct {
	Id        string              `json:"id"`
	Name      string              `json:"name"`
	Reward    ChallengeRewardResp `json:"reward"`
	Status    int                 `json:"status"`
	UpdatedAt int64               `json:"updated_at"`
}

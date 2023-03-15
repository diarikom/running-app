package dto

import "github.com/diarikom/running-app/running-app-api/internal/pkg/ngrule"

type ChallengeRewardResp struct {
	Type     string `json:"type"`
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

type ChallengeRewardClaimResp struct {
	Status int8 `json:"status"`
}

type MilestoneAchievementResp struct {
	CreditReward int64 `json:"credit_reward"`
}

type AchieveChallenge struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Status int         `json:"status"`
	Rules  ngrule.Rule `json:"rules"`
}

package dto

import "time"

type UserChallengeReq struct {
	UserId                  string      `json:"user_id"`
	ChallengeId             string      `json:"-"`
	RewardValue             float64     `json:"-"`
	RewardRefId             string      `json:"-"`
	ChallengeResultSnapshot interface{} `json:"-"`
	Timestamp               time.Time   `json:"-"`
}

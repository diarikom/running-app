package dto

type UserCreditBalanceResp struct {
	Balance         float64 `json:"balance"`
	PendingBalance  float64 `json:"pending_balance"`
	ExpiringBalance float64 `json:"expiring_balance"`
	ExpireTime      int64   `json:"expire_time"`
}

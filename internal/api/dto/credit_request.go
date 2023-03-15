package dto

import (
	"time"
)

type CreditSettleOpt struct {
	TrxId     string
	Notes     string
	ExpiredAt *time.Time
	Timestamp *time.Time
}

type CreditTrxOpt struct {
	Id        string
	WalletId  string
	Amount    float64
	EntryType int8
	Notes     string
	ExpiredAt *time.Time
	Timestamp *time.Time
}

type CreditChargeOpt struct {
	UserId        string
	Amount        float64
	WalletVersion int
}

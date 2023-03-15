package dto

type PageReq struct {
	Skip  int64
	Limit int8
}

type UserResourcesReq struct {
	PageReq
	UserId string
}

package nfacebook

type ErrorResp struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
	TraceId string `json:"fbtrace_id"`
}

type InspectTokenResp struct {
	TokenData
	Error *ErrorResp `json:"error"`
}

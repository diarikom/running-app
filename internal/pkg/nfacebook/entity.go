package nfacebook

type TokenData struct {
	AppId               string   `json:"app_id"`
	Type                string   `json:"type"`
	Application         string   `json:"application"`
	DataAccessExpiresAt int64    `json:"data_access_expires_at"`
	TokenExpiresAt      int64    `json:"expires_at"`
	IsValid             bool     `json:"is_valid"`
	IssuedAt            int64    `json:"issued_at"`
	Scopes              []string `json:"scopes" `
	UserId              string   `json:"user_id"`
}

package models

type TokenResponse struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	AccessExpireTime  int64  `json:"access_expire_time"`
	RefreshExpireTime int64  `json:"refresh_expire_time,omitempty"`
}

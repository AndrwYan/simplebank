package token

import "time"

type Maker interface {

	// CreateToken 创建token的函数
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken 验证token的函数
	VerifyToken(token string) (*Payload, error)
}

package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

//PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker 创建一个新的PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken 创建token
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken 校验token
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	//用于存储解析解析后的payload
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrValidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

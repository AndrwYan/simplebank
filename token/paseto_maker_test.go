package token

import (
	"github.com/AndrewLoveMei/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute

	issueAt := time.Now()
	expiredAt := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//验证token
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Minute)
}

//编写测试token过期的函数
func TestExpiredPasetoToken(t *testing.T) {

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	//生成token
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//测试token过期
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

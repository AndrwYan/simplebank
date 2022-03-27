package api

import (
	"context"
	"github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/AndrewLoveMei/simplebank/db/util"
	"github.com/stretchr/testify/require" //导入Golang的测试工具包
	"testing"
	"time"
)

//创建账户
func createRandomUser(t *testing.T) api.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := api.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.FullName, user1.FullName)
	require.Equal(t, user2.HashedPassword, user1.HashedPassword)

	require.WithinDuration(t, //比较时间
		user2.CreatedAt,
		user1.CreatedAt,
		time.Second,
	)
	require.WithinDuration(t,
		user2.PasswordChangedAt,
		user1.PasswordChangedAt,
		time.Second,
	)

}

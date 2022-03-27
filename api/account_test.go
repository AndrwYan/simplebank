// golang的代码约定就是把测试文件和文件放在同一个文件夹
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
func createAccounts(t *testing.T) api.Account {
	arg := api.CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createAccounts(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createAccounts(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Balance, account1.Balance)
	require.Equal(t, account2.ID, account1.ID)

	require.WithinDuration(
		t,
		account1.CreatedAt,
		account2.CreatedAt,
		time.Second,
	)

}

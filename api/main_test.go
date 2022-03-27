package api

import (
	"database/sql"
	"github.com/AndrewLoveMei/simplebank/db/sqlc"
	db "github.com/AndrewLoveMei/simplebank/db/sqlc"
	"github.com/AndrewLoveMei/simplebank/db/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

var testQueries *api.Queries
var testDB *sql.DB
var err error

func newTestServer(t *testing.T, store *db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("can not load config!", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db!")
	}

	testQueries = api.New(testDB)
	os.Exit(m.Run())
}

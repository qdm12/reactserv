package data

import (
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/reactserv/internal/data/json"
	"github.com/qdm12/reactserv/internal/data/memory"
	"github.com/qdm12/reactserv/internal/data/psql"
	"github.com/qdm12/reactserv/internal/models"
)

type Database interface {
	Close() error
	CreateUser(user models.User) (err error)
	GetUserByID(id uint64) (user models.User, err error)
}

func NewMemory() (Database, error) {
	return memory.NewDatabase()
}

func NewJSON(filepath string) (Database, error) {
	memoryDatabase, err := memory.NewDatabase()
	if err != nil {
		return nil, err
	}
	return json.NewDatabase(memoryDatabase, filepath)
}

func NewPostgres(host, user, password, database string, logger logging.Logger) (Database, error) {
	return psql.NewDatabase(host, user, password, database, logger)
}

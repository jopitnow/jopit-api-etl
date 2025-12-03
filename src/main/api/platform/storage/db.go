package storage

import (
	"os"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
)

func NewNoSQL() *gonosql.Data {
	return gonosql.NewNoSQL(os.Getenv("MONGODB_CONN_STRING"))
}

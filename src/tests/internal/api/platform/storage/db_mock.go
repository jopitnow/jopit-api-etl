package storage

import (
	"context"
	"fmt"
	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	data *gonosql.Data
	once sync.Once
)

func OpenNoSQLMock(server *memongo.Server) *gonosql.Data {
	once.Do(func() {
		initNoSQLMock(server)
	})

	return data
}

func CloseNoSQLMock() {
	if data == nil {
		return
	}

	if err := data.DB.Disconnect(context.Background()); err != nil {
		fmt.Println("Error disconnecting database.", err)
	}

	fmt.Println("Connection closed successfully.")
}

func initNoSQLMock(server *memongo.Server) {
	var (
		errDB    error
		database *mongo.Database
	)
	db, err := getConnection(server)
	if err != nil {
		errDB = fmt.Errorf("error NoSQL connection: %s", err)
	} else {
		// Check the connections
		if err = db.Ping(context.Background(), nil); err != nil {
			errDB = fmt.Errorf("error NoSQL connection: %s", err)
		}
		database = db.Database(memongo.RandomDatabase())
	}

	data = &gonosql.Data{
		DB:       db,
		Error:    errDB,
		Database: database,
	}
}

func getConnection(server *memongo.Server) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(server.URI())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	return mongo.Connect(ctx, clientOptions)
}

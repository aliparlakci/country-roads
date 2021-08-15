package common

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitilizeDb(uri string) (*mongo.Client, func()) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(c, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	close := func() {
		defer cancel()
		if err := client.Disconnect(c); err != nil {
			panic(err)
		}
	}

	return client, close
}

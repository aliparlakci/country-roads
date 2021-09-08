package common

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDb(uri, name, username, password string) (database *mongo.Database, close func()) {
	c, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client()
	if username != "" {
		credentials := options.Credential{
			Username: username,
			Password: password,
		}
		clientOptions = clientOptions.SetAuth(credentials)
	}

	client, err := mongo.Connect(c, clientOptions.ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	database = client.Database(name)

	close = func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}

	return
}

func RedisInitializer(uri, password string) func(db int) *redis.Client {
	return func(db int) *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr: uri,
			DB:   db,
		})
	}
}

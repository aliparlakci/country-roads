package repositories

import (
	"context"
	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactCollection struct {
	Collection *mongo.Collection
}

type ContactRepository interface {
	FindOne(ctx context.Context, filter interface{}) (models.Contact, error)
	InsertOne(ctx context.Context, candidate models.Contact) (interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error
}

type ContactFinder interface {
	FindOne(ctx context.Context, filter interface{}) (models.Contact, error)
}

type ContactInserted interface {
	InsertOne(ctx context.Context, candidate models.Contact) (interface{}, error)
}

type ContactUpdater interface {
	UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error
}


func (c *ContactCollection) FindOne(ctx context.Context, filter interface{}) (models.Contact, error) {
	var contact models.Contact
	result := c.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return contact, err
	}
	err := result.Decode(&contact)
	return contact, err
}

func (c *ContactCollection) InsertOne(ctx context.Context, candidate models.Contact) (interface{}, error) {
	result, err := c.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (c *ContactCollection) UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error {
	_, err := c.Collection.UpdateOne(ctx, filter, changes)
	return err
}

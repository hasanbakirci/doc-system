package mongoClient

import (
	"context"
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ConnectDb(settings config.MongoSettings) (db *mongo.Database, err error) {
	uri := settings.Uri

	log.Infof("Mongo:Connection Uri:%s", uri)
	clientOptions := options.
		Client().
		ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Errorf("Mongo: couldn't connect to mongo: %v", err)
		return db, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Errorf("Mongo: mongo client couldn't connect with background context: %v", err)
		return db, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Errorf("Mongo: Client Ping error", err)
	}

	db = client.Database(settings.DatabaseName)

	return db, err
}

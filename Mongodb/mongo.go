package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const connectionString string = "mongodb+srv://Vishanksingh:7997@cluster0.ga4iiwd.mongodb.net/todo"

const dbName string = "todo"
const collectionName string = "todo"

var Collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		logger.Error("Could not connect to Mongo db at ", zap.String("connectionString", connectionString), zap.Error(err))
		Collection = nil
	}

	logger.Info("Connected to MongoDb")

	Collection = client.Database(dbName).Collection(collectionName)
}

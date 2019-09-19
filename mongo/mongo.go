package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
client.Connect(ctx)
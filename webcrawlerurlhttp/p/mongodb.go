package p

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	StringConnection string
}

func (m *MongoDB) GetConnection() (client *mongo.Client, ctx context.Context) {
	clientOptions := options.Client().ApplyURI(m.StringConnection)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Println(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
	}

	return
}

func (m *MongoDB) Upsert(update interface{}, filter bson.M) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	_, err := collectionMongo.UpdateOne(
		context.Background(),
		filter,
		bson.M{"$set": update},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Println(err)
	}
}

func (m *MongoDB) FindOne(filter bson.M) (result bson.M, err error) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	err = collectionMongo.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		result = bson.M{}
	}

	return
}

func (m *MongoDB) FindAll(filter bson.M) (result []bson.M, err error) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	cursor, err := collectionMongo.Find(ctx, filter)
	if err != nil {
		log.Println(err)
	}

	for cursor.Next(ctx) {
		var resultTemp bson.M
		err := cursor.Decode(&resultTemp)
		if err != nil {
			log.Println(err)
		}

		result = append(result, resultTemp)
	}

	return
}

func (m *MongoDB) DeleteAll(filter bson.M) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	collectionMongo.DeleteMany(context.Background(), filter)
}

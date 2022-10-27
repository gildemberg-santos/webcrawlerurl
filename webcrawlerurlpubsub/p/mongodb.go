package p

import (
	"context"
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
		Logs("Error MongoDB: ", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		Logs("Error MongoDB: ", err)
	}

	return
}

func (m *MongoDB) UpsertOne(update interface{}, filter bson.M) {
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
		Logs("Error MongoDB: ", err)
	}
}

func (m *MongoDB) InsertOne(insert interface{}) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	_, err := collectionMongo.InsertOne(context.Background(), insert)
	if err != nil {
		Logs("Error MongoDB: ", err)
	}
}

func (m *MongoDB) InsertMany(insert bson.A) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	option := options.InsertMany().SetOrdered(false)
	_, err := collectionMongo.InsertMany(context.Background(), insert, option)
	if err != nil {
		Logs("Error MongoDB: ", err)
	}
}

func (m *MongoDB) UpdateOne(update interface{}, filter bson.M) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	_, err := collectionMongo.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		Logs("Error MongoDB: ", err)
	}
}

func (m *MongoDB) UpdateMany(update interface{}, filter bson.M) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")
	_, err := collectionMongo.UpdateMany(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		Logs("Error MongoDB: ", err)
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

func (m *MongoDB) FindAll(filter bson.M, sortResult bson.M) (result []bson.M, err error) {
	client, ctx := m.GetConnection()
	defer client.Disconnect(ctx)

	collectionMongo := client.Database("webcrawler").Collection("links")

	option := options.FindOptions{}
	option.SetProjection(sortResult)
	cursor, err := collectionMongo.Find(ctx, filter, &option)
	if err != nil {
		Logs("Error MongoDB: ", err)
	}

	for cursor.Next(ctx) {
		var resultTemp bson.M
		err := cursor.Decode(&resultTemp)
		if err != nil {
			Logs("Error MongoDB: ", err)
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

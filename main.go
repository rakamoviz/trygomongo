package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

func upsertEvent(ctx context.Context, collection *mongo.Collection, event *Event) error {
	filter := bson.M{"_id": event.ID}
	update := bson.M{"$set": event}
	upsert := true

	_, err := collection.UpdateOne(ctx, filter, update, &options.UpdateOptions{
		Upsert: &upsert,
	})

	return err
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)

	time.Sleep(1 * time.Second)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	eventsCollection := client.Database("trygomongo").Collection("events")

	objectID := primitive.NewObjectID()

	event := &Event{
		ID:        objectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Text:      "Example event",
		Completed: false,
	}

	err = upsertEvent(ctx, eventsCollection, event)
	if err != nil {
		log.Fatal(err)
	}

	event2 := &Event{
		ID:        objectID,
		CreatedAt: event.CreatedAt,
		UpdatedAt: time.Now(),
		Text:      event.Text,
		Completed: true,
	}

	err = upsertEvent(ctx, eventsCollection, event2)
	if err != nil {
		log.Fatal(err)
	}
}

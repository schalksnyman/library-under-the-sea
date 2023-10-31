package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	repo "library-under-the-sea/services/library-repo/domain"
	library "library-under-the-sea/services/library/domain"
	"log"
)

func NewMongoClient(connectString string, dbName string) repo.DBHandler {
	clientOptions := options.Client().ApplyURI(connectString)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return &client{}
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return &client{}
	}

	return &client{
		MongoClient: *mongoClient,
		database:    mongoClient.Database(dbName),
	}
}

type client struct {
	MongoClient mongo.Client
	database    *mongo.Database
}

func (c *client) Get(ctx context.Context, id string) (*library.Book, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result := c.database.Collection("books").FindOne(context.Background(), bson.M{"_id": objID})
	var book library.Book

	err = result.Decode(book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (c *client) ListByTitle(ctx context.Context, title string) ([]*library.Book, error) {
	var results []*library.Book
	collection := c.database.Collection("books")
	filter := bson.D{{"title", title}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem library.Book
		err2 := cur.Decode(&elem)
		if err2 != nil {
			log.Fatal(err2)
		}
		results = append(results, &elem)
	}
	return results, nil
}

func (c *client) ListAll(ctx context.Context) ([]*library.Book, error) {
	var results []*library.Book
	collection := c.database.Collection("books")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem library.Book
		err2 := cur.Decode(&elem)
		if err2 != nil {
			log.Fatal(err2)
		}
		results = append(results, &elem)
	}
	return results, nil
}

func (c *client) Save(ctx context.Context, book library.Book) (string, error) {
	collection := c.database.Collection("books")
	res, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		return "", err
	}

	objID := res.InsertedID.(primitive.ObjectID).Hex()
	return objID, nil
}

func (c *client) Delete(ctx context.Context, id string) error {
	collection := c.database.Collection("books")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		return err
	}

	if res.DeletedCount != 1 {
		return errors.New("deleted count is not one")
	}

	return nil
}

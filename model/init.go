package model

import (
	"context"

	"github.com/asynccnu/classroom_service_v2/log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Self *mongo.Client
}

var DB *Database

func InitSelfDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(viper.GetString("db.url"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		// log.Errorf(err, "Database connection failed.")
		log.Error("Database connection failed.")
		return client
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("Database connection failed: " + err.Error())
		panic(err)
	}

	log.Info("Connected to MongoDB!")

	return client
}

func GetSelfDB() *mongo.Client {
	return InitSelfDB()
}

func (db *Database) Init() {
	DB = &Database{
		Self: GetSelfDB(),
	}

	DBName = viper.GetString("db.name")
}

func (db *Database) Close() {
	DB.Self.Disconnect(context.TODO())
}

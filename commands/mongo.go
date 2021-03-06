package commands

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)


var mongodbSession *mgo.Session

func init() {
	RootCmd.PersistentFlags().String("mongodb_uri", "localhost", "host where mongoDB is")
	viper.BindPFlag("mongodb_uri", RootCmd.PersistentFlags().Lookup("mongodb_uri"))
	viper.SetDefault("dbname", "planet")
	CreateUniqueIndexes()
}

func DBSession() *mgo.Session {
	if mongodbSession == nil {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			uri = viper.GetString("mongodb_uri")

			if uri == "" {
				log.Fatalln("No connection uri for MongoDB provided")
			}
		}

		var err error
		mongodbSession, err = mgo.Dial(uri)
		if mongodbSession == nil || err != nil {
			log.Fatalf("Can't connect to mongo, go error %v\n", err)
		}

		mongodbSession.SetSafe(&mgo.Safe{})
	}
	return mongodbSession
}

func Items() *mgo.Collection {
	return DB().C("items")
}

func Channels() *mgo.Collection {
	return DB().C("channels")
}

func DB() *mgo.Database {
	return DBSession().DB(viper.GetString("dbname"))
}

func CreateUniqueIndexes() {
	idx := mgo.Index{
		Key:        []string{"key"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := Items().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}

	if err := Channels().EnsureIndex(idx); err != nil {
		fmt.Println(err)
	}
}
package db

import (
	goHome "github.com/mitchellh/go-homedir"
	bolt "go.etcd.io/bbolt"
)

type DbClient struct {
	Db *bolt.DB
}

var dbClientSingletionInstance *DbClient = nil

// NewConnection opens a new DB connection and returns with a new DbClient instance as result.
func NewConnection() DbClient {
	if dbClientSingletionInstance == nil {
		homeDir, err := goHome.Dir()
		if err != nil {
			panic("error occourred getting the home directory")
		}

		db, err := bolt.Open(homeDir+"/bolt_db/TODO.db", 0666, nil)
		if err != nil {
			err.Error()
			panic(err)
		}
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("todo_new"))
			return err
		})

		if err != nil {
			err.Error()
			panic(err)
		}
		dbClientSingletionInstance = &DbClient{db}
	}
	return *dbClientSingletionInstance
}

func (dbClient *DbClient) CloseConnection() {
	dbClient.Db.Close()
	dbClientSingletionInstance = nil
}

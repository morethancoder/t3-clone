package db

import (
	"errors"
	"fmt"
	"morethancoder/t3-clone/utils"
	"net/http"
	"os"
)

type db struct {
	Url string
	Client *http.Client
}

var Db *db

func init() {
	var err error
	Db, err = NewDB()
	if err != nil {
		utils.Log.Fatal(err.Error())
	}
}

func NewDB() (*db, error) {
	if os.Getenv("DB_URL") == "" {
		return nil, errors.New("[ERROR] env DB_URL is not set")
	}
	return &db{
		Url: os.Getenv("DB_URL"),
		Client: &http.Client{},
	}, nil
}

func FileUrl(collectionId, recordId, filename string) string {
	return fmt.Sprintf("%s/api/files/%s/%s/%s?token=", Db.Url, collectionId ,recordId ,filename)
}





package migrations

import (
  "manga-app/db/dbClients"
  "github.com/golang/glog"
)

func Run () {
  client := dbClients.NewDynamoDbClient()

  err := createMangaTable(client)
  if err != nil {
    glog.Error(err)
  }
}

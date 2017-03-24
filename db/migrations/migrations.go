package migrations

import (
  dbClients "manga-app/db/clients"
  "github.com/golang/glog"
)

func Run () {
  client := dbClients.NewDynamoDbClient()

  err := createMangaTable(client)
  if err != nil {
    glog.Error(err)
    panic(err)
  }
}

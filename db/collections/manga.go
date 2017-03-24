package collections

import (
  "fmt"

  "manga-app/models"
  "manga-app/db/collections/tableNames"

  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MangaCollection struct {
  Collection // embedded
}


/**
 * Returns instance of the manga dao
 */
func GetMangaCollection() (*MangaCollection, error) {
  var mangaCollection *MangaCollection = nil
  embeddedCollection, err := initCollection(tableNames.MANGA_TABLE)
  if err != nil {
    return mangaCollection, err
  }

  mangaCollection = &MangaCollection{embeddedCollection}

  return mangaCollection, nil
}

func (m MangaCollection) BatchPutItems (manga []models.Manga) error {
  // type cast manga to that of a regular interface
  items := make([]map[string]*dynamodb.AttributeValue, len(manga))
  for i, value := range manga {
    item, err := dynamodbattribute.MarshalMap(&value)
    if err != nil {
      fmt.Println(err)
      continue
    }

    items[i] = item
  }

  return m.Collection.BatchPutItems(items)
}


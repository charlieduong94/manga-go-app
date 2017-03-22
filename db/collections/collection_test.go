package collections

import (
  "flag"
  "testing"
  "manga-app/config"
  "manga-app/models"
  "manga-app/db/dbClients"
  "manga-app/db/migrations"
  "manga-app/db/collections/tableNames"

  "github.com/satori/go.uuid"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/stretchr/testify/assert"
)

func TestMain (m *testing.M) {
  flag.Parse()
  config.Load("test/config.yml")
  migrations.Run()
  m.Run()
}

func TestBatchPutItems (t *testing.T) {
  mangaList := models.MangaList{}

  mangaList.Manga = make([]models.Manga, 1)

  uniqueIdA := uuid.NewV4().String()

  originalManga := models.Manga{
    Id: uniqueIdA,
    Title: "mangaA",
    Image: "99/9949c70030a89c9a2a1d5273a627de77ac2aaa948c961f1212c2ba46.jpg",
    Alias: "mangaA",
    Status: 1,
    Category: []string{
      "Shounen",
    },
    LastChapterDate: 0,
    Hits: 0,
  }

  mangaList.Manga[0] = originalManga

  collection, _ := GetMangaCollection()

  _ = collection.BatchPutItems(mangaList.Manga)

  client := dbClients.NewDynamoDbClient()

  queryInput := &dynamodb.QueryInput{
    TableName: aws.String(tableNames.MANGA_TABLE),
    KeyConditionExpression: aws.String("id = :i"),
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":i": {
        S: aws.String(uniqueIdA),
      },
    },
  }

  output, _ := client.Query(queryInput)


  // marshal the data back into a Manga obj
  retrievedManga := &models.Manga{}
  _ = dynamodbattribute.UnmarshalMap(output.Items[0], retrievedManga)

  assert := assert.New(t)
  assert.Equal(originalManga.Id, retrievedManga.Id)
  assert.Equal(originalManga.Title, retrievedManga.Title)
  assert.Equal(originalManga.Image, retrievedManga.Image)
  assert.Equal(originalManga.Alias, retrievedManga.Alias)
  assert.Equal(originalManga.Status, retrievedManga.Status)
  assert.Equal(originalManga.Category, retrievedManga.Category)
  assert.Equal(originalManga.LastChapterDate, retrievedManga.LastChapterDate)
  assert.Equal(originalManga.Hits, retrievedManga.Hits)
}
/**
 * Performs api calls to manga eden
 */
package manga

import (
  "github.com/golang/glog"
  "log"
  "sync"
  "manga-app/db/collections"
  "manga-app/models"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type MangaService struct {
  collection *collections.MangaCollection
}

const LIST_URL = "https://www.mangaeden.com/api/list/0/"

var instance *MangaService = nil

var once sync.Once

func GetInstance () *MangaService {
  once.Do(func () {
    glog.Info("Instantiating manga service")
    collection, err := collections.GetMangaCollection()
    if err != nil {
      log.Fatal(err)
    }
    instance = &MangaService{collection}
  })

  return instance
}

func listManga () (models.MangaList, error) {
  var pagedListInput models.MangaListInput
  var mangaList models.MangaList

  res, err := http.Get(LIST_URL)
  defer res.Body.Close()
  if (err != nil) {
    return mangaList, err
  }

  data, err := ioutil.ReadAll(res.Body)
  if (err != nil) {
    return mangaList, err
  }

  err = json.Unmarshal(data, &pagedListInput)
  if (err != nil) {
    return mangaList, err
  }

  preTransformedMangaList := pagedListInput.Manga

  length := len(preTransformedMangaList)

  manga := make([]models.Manga, length)

  // transform the data into a format that we like
  for i := 0; i < length; i++ {
    curManga := preTransformedMangaList[i]
    manga[i] = models.Manga{
      "english", // TODO: add support for others later
      curManga.Id,
      curManga.Title,
      curManga.Image,
      curManga.Alias,
      curManga.Status,
      curManga.Category,
      curManga.LastChapterDate,
      curManga.Hits,
    }
  }

  mangaList = models.MangaList{manga}

  return mangaList, nil
}

func (m MangaService) SyncManga () error {
  mangaList, err := listManga()
  if err != nil {
    return err
  }

  err = m.collection.BatchPutItems(mangaList.Manga)
  if err != nil {
    return err
  }

  return nil
}

func (m MangaService) GetLatestUpdates () ([]models.Manga, error) {
  query := &dynamodb.QueryInput {
    IndexName: aws.String("lastChapterIndex"),
    Limit: aws.Int64(25),
    KeyConditionExpression: aws.String("lang = :l and lastChapterDate > :i"),
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
      ":l": {
        S: aws.String("english"),
      },
      ":i": {
        N: aws.String("0"),
      },
    },
  }

  items, err := m.collection.Query(query)
  if err != nil {
    return make([]models.Manga, 0), err
  }

  return items, nil
}

/*
func GetMostPopular () models.MangaList {
*/

/**
 * Performs api calls to manga eden
 */
package manga

import (
  "errors"
  "encoding/base64"
  "strings"
  "log"
  "sync"
  "net/http"
  "encoding/json"
  "io/ioutil"

  "manga-app/db/collections"
  "manga-app/models"

  "github.com/golang/glog"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func encodeKey (key map[string]*dynamodb.AttributeValue) (string, error) {
  keyMap := make(map[string]string)

  err := dynamodbattribute.UnmarshalMap(key, &keyMap)
  byteArray := []byte(keyMap["id"] + "|" + keyMap["lastChapterDate"] + "|" + keyMap["lang"])
  if err != nil {
    return "", err
  }

  return base64.StdEncoding.EncodeToString(byteArray), nil
}

func decodeKey (str string) (map[string]*dynamodb.AttributeValue, error) {
  decodedMap := make(map[string]*dynamodb.AttributeValue)

  decodedStr, err := base64.StdEncoding.DecodeString(str)
  if err != nil {
    return decodedMap, err
  }

  parts := strings.Split(string(decodedStr), "|")
  if len(parts) != 3 {
    return decodedMap, errors.New("Unable to decode last key")
  }

  glog.Info("parts", parts)

  decodedMap["id"] = &dynamodb.AttributeValue{ S: aws.String(parts[0]) }
  decodedMap["lastChapterDate"] = &dynamodb.AttributeValue{ N: aws.String(parts[1]) }
  decodedMap["lang"] = &dynamodb.AttributeValue{ S: aws.String(parts[2]) }

  return decodedMap, nil
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

  mangaList.Manga = manga

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

func (m MangaService) GetLatestUpdates (startKey string) (models.MangaList, error) {
  query := &dynamodb.QueryInput {
    IndexName: aws.String("lastChapterIndex"),
    Limit: aws.Int64(25),
    KeyConditionExpression: aws.String("lang = :l and lastChapterDate > :i"),
    ScanIndexForward: aws.Bool(false),
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
      ":l": {
        S: aws.String("english"),
      },
      ":i": {
        N: aws.String("0"),
      },
    },
  }

  var mangaList models.MangaList

  if len(startKey) > 0 {
    newKey, err := decodeKey(startKey)
    if err != nil {
      return mangaList, errors.New("Invalid startKey provided")
    }
    glog.Info(newKey)
    query.ExclusiveStartKey = newKey
  }


  output, err := m.collection.Query(query)
  items := output.Items
  manga := make([]models.Manga, len(items))

  for i, v := range items {
    var item models.Manga

    err := dynamodbattribute.ConvertFromMap(v, &item)
    if err != nil {
      return mangaList, errors.New("Problem converting data")
    }

    manga[i] = item
  }

  lastKey, err := encodeKey(output.LastEvaluatedKey)
  if err != nil {
    return mangaList, err
  }

  mangaList.Manga = manga
  mangaList.LastKey = &lastKey

  if err != nil {
    return mangaList, err
  }

  return mangaList, nil
}

/*
func GetMostPopular () models.MangaList {
*/

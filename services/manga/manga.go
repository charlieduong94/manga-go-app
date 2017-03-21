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

/*
func GetLatestUpdates () models.MangaList {

}

func GetMostPopular () models.MangaList {
*/

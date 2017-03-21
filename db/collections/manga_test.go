package collections

import (
  "testing"
  "manga-app/models"
)

func TestBatchPutItems (t *testing.T) {
  mangaList := models.MangaList{}

  mangaList.Manga = make([]models.Manga, 1)

  manga := &models.Manga{}

  mangaList.Manga[0] = *manga

  collection, _ := GetMangaCollection()
  collection.BatchPutItems(mangaList)
}

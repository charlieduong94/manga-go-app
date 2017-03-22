package jobs

import (
  "time"
  "github.com/golang/glog"
  "manga-app/services/manga"
)

var HALF_DAY_IN_SECONDS = 43200 * time.Second

func startMangaSyncJob () {
  // run forever
  mangaService := manga.GetInstance()
  for {
    glog.Info("Syncing with MangaEden...")

    mangaService.SyncManga()

    glog.Info("Manga sync Complete!")
    time.Sleep(HALF_DAY_IN_SECONDS)
  }
}


package jobs

import (
  "time"
  "github.com/golang/glog"
  _ "manga-app/services/manga"
)

var HALF_DAY_IN_SECONDS = 43200 * time.Second

func StartMangaSyncJob () {
  // run forever
  for {
    glog.Info("Syncing with MangaEden...")

    glog.Info("Manga sync Complete!")
    time.Sleep(HALF_DAY_IN_SECONDS)
  }
}


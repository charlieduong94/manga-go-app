package routes

import (
  "github.com/golang/glog"
  "gopkg.in/gin-gonic/gin.v1"
  "manga-app/services/manga"
)

// use service to handle calls
func handleLatestUpdates (m *manga.MangaService) gin.HandlerFunc {
  return func (context *gin.Context) {
    startKey := context.Query("startKey")

    manga, err := m.GetLatestUpdates(startKey)
    if err != nil {
      glog.Error(err)

      context.JSON(500, gin.H{
        "error": err,
      })
    } else {
      context.JSON(200, manga)
    }
  }
}

func handleMostPopular (context *gin.Context) {
  context.JSON(200, gin.H{
    "message": "most popular",
  })
}

func ApplyMangaRoutes (router *gin.Engine) {
  mangaService := manga.GetInstance()

  v1 := router.Group("/v1/manga")
  {
    v1.GET("/latest-updates", handleLatestUpdates(mangaService))
    v1.GET("/most-popular", handleMostPopular)
  }
}

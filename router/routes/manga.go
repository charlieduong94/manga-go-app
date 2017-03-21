package routes

import (
  "gopkg.in/gin-gonic/gin.v1"
  "manga-app/services/manga"
)

var mangaService *manga.MangaService = nil

// use service to handle calls
func handleLatestUpdates (context *gin.Context) {
  context.JSON(200, gin.H{
    "message": "most popular",
  })
}

func handleMostPopular (context *gin.Context) {
  context.JSON(200, gin.H{
    "message": "most popular",
  })
}

func ApplyMangaRoutes (router *gin.Engine) {
  mangaService = manga.GetInstance()

  v1 := router.Group("/v1/manga")
  {
    v1.GET("/latest-updates", handleLatestUpdates)
    v1.GET("/most-popular", handleMostPopular)
  }
}

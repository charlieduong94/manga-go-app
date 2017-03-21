package router

import (
  "gopkg.in/gin-gonic/gin.v1"
  "manga-app/router/routes"
)

/**
 * Loads routes into a http router
 */
func Load () *gin.Engine {
  router := gin.Default()

  routes.ApplyMangaRoutes(router)

  return router
}

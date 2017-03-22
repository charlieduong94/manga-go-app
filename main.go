package main

import (
  "flag"
  "manga-app/config"
  "manga-app/router"
  "manga-app/jobs"
  "manga-app/db/migrations"
)

func main () {
  // parse cmd line args
  flag.Parse()

  config.Load("config.yml")

  migrations.Run()

  jobs.StartAll()

  app := router.Load()
  app.Run()
}

package main

import (
  "flag"
  "manga-app/router"
  "manga-app/jobs"
)

func main () {
  // parse cmd line args
  flag.Parse()

  jobs.StartAll()

  app := router.Load()
  app.Run()
}

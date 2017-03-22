package jobs

/**
 * Start up any async tasks here
 */
func StartAll () {
  go startMangaSyncJob()
}

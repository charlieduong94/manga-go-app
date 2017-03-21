package models

type PagedMangaListInput struct {
  MangaListInput
  Page int
  Start int
  End int
  Total int
}

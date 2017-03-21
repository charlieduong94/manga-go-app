package models

type MangaListInput struct {
  Manga []MangaInput `json:"manga"`
}

type MangaList struct {
  Manga []Manga `json:"manga"`
}

package models

// input
type MangaListInput struct {
  Manga []MangaInput `json:"manga"`
}

// output
type MangaList struct {
  Manga []Manga `json:"manga"`
  LastKey *string `json:"lastKey"`
}

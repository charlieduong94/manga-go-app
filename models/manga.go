package models

type MangaInput struct {
  Id string `json:"i"`
  Title string `json:"t"`
  Image string `json:"im"`
  Alias string `json:"a"`
  Status int `json:"s"`
  Category []string `json:"c"`
  LastChapterDate float64 `json:"ld"`
  Hits int `json:"h"`
}

type Manga struct {
  Language string `json:"lang"`
  Id string `json:"id"`
  Title string `json:"title"`
  Image string `json:"image"`
  Alias string `json:"alias"`
  Status int `json:"status"`
  Category []string `json:"category"`
  LastChapterDate float64 `json:"lastChapterDate"`
  Hits int `json:"hits"`
}

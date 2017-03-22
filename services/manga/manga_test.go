package manga

import (
  "strings"
  "io/ioutil"
  "encoding/json"
  "testing"
  "net/http"
  "github.com/jarcoal/httpmock"
  "manga-app/models"
)

func check (err error, t *testing.T) {
  if err != nil {
    t.Fail()
  }
}

func TestListManga (t *testing.T) {
  httpmock.Activate()
  defer httpmock.DeactivateAndReset()

  expectedData, err := ioutil.ReadFile("test/expected_listmanga_value.json")
  check(err, t)

  var expectedList models.MangaList

  err = json.Unmarshal(expectedData, &expectedList)
  check(err, t)

  httpmock.RegisterResponder("GET", LIST_URL,
    func (req *http.Request) (*http.Response, error) {
      data, err := ioutil.ReadFile("test/ListManga_http_response.json")
      check(err, t)

      var returnedData models.MangaListInput

      err = json.Unmarshal(data, &returnedData)
      check(err, t)

      response, err := httpmock.NewJsonResponse(200, returnedData)
      check(err, t)
      return response, nil
    },
  )

  actualList, err := listManga()
  check(err, t)

  // marshal objects into string for assertion
  actual, _ := json.Marshal(actualList)
  expected, _ := json.Marshal(expectedList)

  actualJson := string(actual)
  expectedJson := string(expected)

  if strings.Compare(actualJson, expectedJson) != 0 {
    t.Log("Actual and expected results are not equal")
    t.Log("Actual:", actualJson)
    t.Log("Expected:", expectedJson)
    t.Fail()
  }
}

package manga

import (
  "strings"
  "io/ioutil"
  "encoding/json"
  "encoding/base64"
  "testing"
  "net/http"
  "github.com/jarcoal/httpmock"
  "manga-app/models"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/stretchr/testify/assert"
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

  for i := range expectedList.Manga {
    expectedList.Manga[i].Language = "english"
  }

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
  actual, _ := json.Marshal(actualList.Manga)
  expected, _ := json.Marshal(expectedList.Manga)

  actualJson := string(actual)
  expectedJson := string(expected)

  if strings.Compare(actualJson, expectedJson) != 0 {
    t.Log("Actual and expected results are not equal")
    t.Log("Actual:", actualJson, len(actualJson))
    t.Log("Expected:", expectedJson, len(expectedJson))
    t.Fail()
  }
}

func TestEncodeKey (t *testing.T) {
  testId := "id"
  testLastChapterDate := "1000000"
  testLang := "english"

  testMap := make(map[string]*dynamodb.AttributeValue)
  testMap["id"] = &dynamodb.AttributeValue{ S: aws.String(testId) }
  testMap["lastChapterDate"] = &dynamodb.AttributeValue{ N: aws.String(testLastChapterDate) }
  testMap["lang"] = &dynamodb.AttributeValue{ S: aws.String(testLang) }

  newKey, err:= encodeKey(testMap)
  if err != nil {
    t.Log("Failed to encode map")
    t.Fail()
  }

  byteArray, err := base64.StdEncoding.DecodeString(newKey)
  if err != nil {
    t.Log("Failed to decode base64 string")
    t.Fail()
  }

  assert := assert.New(t)
  assert.Equal(string(byteArray), testId + "|" + testLastChapterDate + "|" + testLang)
}

func TestDecodeKey (t *testing.T) {
  encoded := "aWR8MTAwMDAwMHxlbmdsaXNo"
  decoded := "id|1000000|english"
  key, err := decodeKey(encoded)
  if err != nil {
    t.Log("Failed to decode string")
    t.Fail()
  }
  assert := assert.New(t)
  assert.Equal(*key["id"].S + "|" + *key["lastChapterDate"].N + "|" + *key["lang"].S, decoded)
}

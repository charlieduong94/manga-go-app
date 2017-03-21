package collections

import (
  "fmt"
  "github.com/golang/glog"
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "manga-app/db/dbClients"
  "manga-app/models"
)

type MangaCollection struct {
  dbClient *dynamodb.DynamoDB
}

const TABLE_NAME = "golang-manga-test"

func createTable (client *dynamodb.DynamoDB) error {
  glog.Info("Initializing table...")

  tableName := aws.String(TABLE_NAME)
  createParams := &dynamodb.CreateTableInput{
    TableName: tableName,
    AttributeDefinitions: []*dynamodb.AttributeDefinition{
      {
        AttributeName: aws.String("id"),
        AttributeType: aws.String("S"),
      },
    },
    KeySchema: []*dynamodb.KeySchemaElement{
      {
        AttributeName: aws.String("id"),
        KeyType: aws.String("HASH"),
      },
    },
    ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
      ReadCapacityUnits: aws.Int64(10),
      WriteCapacityUnits: aws.Int64(10),
    },
  }

  _, err := client.CreateTable(createParams)
  if err != nil {
    // check AND cast err if typeof awserr.Error
    if awsErr, ok := err.(awserr.Error); ok {
      // handle codes accordingly
      code := awsErr.Code()
      if (code == "ResourceInUseException") {
        glog.Info("Table already created")
      }
    } else {
      // return the error that cannot be handled
      return err
    }
  }

  describeParams := &dynamodb.DescribeTableInput {
    TableName: tableName,
  }

  ready := false
  for ready != true {
    output, err := client.DescribeTable(describeParams)
    if err != nil {
      return err
    }

    // dereference table status for comparison
    tableStatus := *output.Table.TableStatus

    if tableStatus == "ACTIVE" {
      ready = true
    } else {
      // give the table time to complete creation
      time.Sleep(5 * time.Second)
    }
  }

  glog.Infof("Table %s is ready", TABLE_NAME)
  return nil
}

/**
 * Returns instance of the manga dao
 */
func GetMangaCollection() (*MangaCollection, error) {
  glog.Info("Creating MangaCollection")
  var collection *MangaCollection
  client := dbClients.NewDynamoDBClient()

  err := createTable(client)
  if err != nil {
    glog.Error(err)
    return collection, err
  }

  collection = &MangaCollection{client}

  return collection, nil
}

func (collection *MangaCollection) BatchPutItems (mangaList models.MangaList) error {
  length := len(mangaList.Manga)
  writeRequests := make([]*dynamodb.WriteRequest, length)

  // build write requests
  for i := 0; i < length; i++ {
    writeRequest := &dynamodb.WriteRequest{}
    item, err := dynamodbattribute.MarshalMap(&mangaList.Manga[i])
    glog.Info(item)
    fmt.Print(item)
    if err != nil {
      glog.Error(err)
      continue
    }

    putRequest := &dynamodb.PutRequest{}
    putRequest.Item = item
    writeRequest.SetPutRequest(putRequest)
    writeRequests[i] = writeRequest
  }

  glog.Info(writeRequests)
  return nil
}
/*

func (collection *MangaCollection) GetLatestUpdates () models.MangaList {
  // TODO: implement
}

func (collection *MangaCollection) GetMostPopular () models.MangaList {
  // TODO: implement
}
*/

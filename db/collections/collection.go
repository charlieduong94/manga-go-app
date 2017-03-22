package collections

import (
  "errors"
  "github.com/golang/glog"
  "manga-app/db/dbClients"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type Collection struct {
  dbClient *dynamodb.DynamoDB
  TableName string
}

func initCollection (tableName string) (Collection, error) {
  var collection Collection
  if len(tableName) == 0 {
    return collection, errors.New("A tableName must be provided")
  }

  client := dbClients.NewDynamoDbClient()

  collection.dbClient = client
  collection.TableName = tableName

  return collection, nil
}

func (collection *Collection) BatchPutItems (items []map[string]*dynamodb.AttributeValue) error {
  client := collection.dbClient
  length := len(items)
  writeRequests := make([]*dynamodb.WriteRequest, length)

  // build write requests
  for i := 0; i < length; i++ {
    writeRequest := &dynamodb.WriteRequest{}

    putRequest := &dynamodb.PutRequest{}
    putRequest.Item = items[i]

    writeRequest.SetPutRequest(putRequest)
    writeRequests[i] = writeRequest
  }

  tableWrites := map[string][]*dynamodb.WriteRequest{}
  tableWrites[collection.TableName] = writeRequests

  batchWriteInput := &dynamodb.BatchWriteItemInput{
    RequestItems: tableWrites,
  }

  output, err := client.BatchWriteItem(batchWriteInput)
  if err != nil {
    glog.Info("err", err)
    return err
  }

  glog.Info("Output", *output)

  return nil
}

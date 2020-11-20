package main

import (
	"errors"
	"math/rand"
	"syscall"
	"time"

	"github.com/guregu/dynamo"
	"github.com/slack-go/slack"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-lambda-go/lambda"
)

type Config struct {
	slack    Slack
	dynamodb DynamoDB
}

type Slack struct {
	webHookUrl string
}

type DynamoDB struct {
	region    string
	tableName string
}

type Word struct {
	WordID int    `json: "WordId"`
	Word   string `json: "Word"`
}

func newSlackBot() (Config, error) {
	var config Config
	webHookUrl, found := syscall.Getenv("WEB_HOOK_URL")
	if !found {
		return config, errors.New("Web hook URL is not found")
	}
	config.slack.webHookUrl = webHookUrl
	dynamoRegion, found := syscall.Getenv("DYNAMO_REGION")
	if !found {
		return config, errors.New("DynamoDB region is not found")
	}
	config.dynamodb.region = dynamoRegion
	dynamoTableName, found := syscall.Getenv("DYNAMO_TABLE_NAME")
	if !found {
		return config, errors.New("DynamoDB table name is not found")
	}
	config.dynamodb.tableName = dynamoTableName

	return config, nil
}

func (c *Config) getRandomWordFromDynamoDB() (result *Word, err error) {
	svc := dynamo.New(session.New(), &aws.Config{Region: aws.String(c.dynamodb.region)})
	table := svc.Table(c.dynamodb.tableName)

	tableDesc, err := table.Describe().Run()
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Int63n(tableDesc.Items)

	err = table.Get("WordId", randomNum).One(&result)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Config) postMessageToSlack(result *Word) error {

	webHookMsg := slack.WebhookMessage{
		Username: "Iron Man",
		Text:     result.Word,
	}
	err := slack.PostWebhook(c.slack.webHookUrl, &webHookMsg)
	return err
}

func lambdaHandler() error {
	config, err := newSlackBot()

	// Get word from dynamoDB
	result, err := config.getRandomWordFromDynamoDB()
	if err != nil {
		return err
	}
	err = config.postMessageToSlack(result)

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}

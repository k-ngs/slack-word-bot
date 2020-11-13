package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WordByMongo struct {
	Id   string
	Word string
}

type Config struct {
	Slack SlackConfig `json: "slack"`
	Mongo MongoConfig `json: "mongo"`
}

type SlackConfig struct {
	WebHookURL       string `json: "webHookUrl"`
	OAuthAccessToken string `json: "oauthAccessToken"`
	ChannelID        string `json: "channelID"`
}

type MongoConfig struct {
	Host           string `json: "host"`
	Port           string `json: "port"`
	DbName         string `json: "dbName"`
	CollectionName string `json: "collectionName"`
}

// TODO: Test -> Random pick a word from MongoDB
func (config *Config) getRandomWordByMongoDB() (word *WordByMongo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoHost := fmt.Sprintf("mongodb://%s:%s", config.Mongo.Host, config.Mongo.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	col := client.Database(config.Mongo.DbName).Collection(config.Mongo.CollectionName)

	err = col.FindOne(context.Background(), bson.M{"_id": ""}).Decode(&word)

	return
}

func (config *Config) handleSendPayload(msg string) error {
	actionName := "word_value"
	attachment := slack.Attachment{
		Color: "#ff8c00",
		Fields: []slack.AttachmentField{
			{
				Title: "Do you like it?",
			},
		},
		Actions: []slack.AttachmentAction{
			{
				Name:  actionName,
				Text:  "さいこー!",
				Type:  "button",
				Style: "primary",
				Value: "saiko",
			},
			{
				Name:  actionName,
				Text:  "いいね",
				Type:  "button",
				Style: "primary",
				Value: "iine",
			},
			{
				Name:  actionName,
				Text:  "ふつう",
				Type:  "button",
				Style: "default",
				Value: "futsu",
			},
			{
				Name:  actionName,
				Text:  "びみょー",
				Type:  "button",
				Style: "danger",
				Value: "bimyo",
			},
		},
	}

	var attachments []slack.Attachment
	attachments = append(attachments, attachment)

	webHookMsg := slack.WebhookMessage{
		Username:    "PostedNagase",
		Text:        msg,
		Attachments: attachments,
	}

	err := slack.PostWebhook(config.Slack.WebHookURL, &webHookMsg)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Read config file and unmarshal
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	// TODO: Pick the word from MongoDB(use getRandomWordByMongoDB())
	msg := "Test message"

	// Post message to slack
	err = config.handleSendPayload(msg)
	if err != nil {
		panic(err)
	}
}

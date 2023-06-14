package app

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jaswdr/faker"

	"github.com/Eviljeks/test-twitter-feed/pkg/client"
)

const wordsMin = 10
const wordsMax = 15

type Bot struct {
	delaySec       uint
	requestsPerMin uint
	apiClient      *client.ApiClient
	faker          *faker.Faker
}

func NewBot(delaySec uint, requestsPerMin uint, apiClient *client.ApiClient, faker *faker.Faker) *Bot {
	return &Bot{
		delaySec:       delaySec,
		requestsPerMin: requestsPerMin,
		apiClient:      apiClient,
		faker:          faker,
	}
}

func (b *Bot) Run(ctx context.Context) {
	time.Sleep(time.Second * time.Duration(b.delaySec))

	ticker := time.NewTicker(time.Second * time.Duration(60/b.requestsPerMin))

	for {
		select {
		case <-ticker.C:
			content := faker.New().Lorem().Sentence(rand.Intn(wordsMax-wordsMin) + wordsMin)
			err := b.apiClient.SaveMessage(content)
			if err != nil {
				log.Default().Println(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}

}

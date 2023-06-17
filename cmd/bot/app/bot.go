package app

import (
	"context"
	"math/rand"
	"time"

	"github.com/jaswdr/faker"
	"github.com/sirupsen/logrus"

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
	tick := time.Second * time.Duration(60/b.requestsPerMin)

	time.Sleep(time.Second*time.Duration(b.delaySec) - tick)

	ticker := time.NewTicker(tick)

	for {
		select {
		case <-ticker.C:
			content := faker.New().Lorem().Sentence(rand.Intn(wordsMax-wordsMin) + wordsMin)
			err := b.apiClient.SaveMessage(content)
			if err != nil {
				logrus.Errorln(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}

}

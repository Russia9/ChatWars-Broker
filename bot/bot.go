package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"kod-guildbot/bot/messages"
	"kod-guildbot/lib"
	"strconv"
	"time"
)

func InitBot(telegramToken string, logger *logrus.Logger, consumer *kafka.Consumer, db *sql.DB) error {
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:  telegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
	if err != nil {
		return err
	}

	bot.Handle("/start", func(message *telebot.Message) {
		fmt.Println(message.Chat.ID)
	})

	consumer.SubscribeTopics([]string{"cw3-offers"}, nil)

	defer bot.Start()

	chat, err := bot.ChatByID(lib.GetEnvOrDefault("CWBR_CHANNEL_ID", "-1001483067163"))

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message messages.OfferMessage
			err = json.Unmarshal([]byte(msg.Value), &message)
			if err != nil {
				logger.Error(fmt.Sprintf("Decoder error: %v (%v)\n", err, msg))
			}
			msgString :=
				" " + message.SellerCastle + " <code>" + message.SellerName + "</code> : \n" +
				" " + strconv.Itoa(message.Quantity) + " " + message.Item + " * 💰" + strconv.Itoa(message.Price)
			_, err = bot.Send(chat, msgString, telebot.ParseMode(telebot.ModeHTML))
			if err != nil {
				logger.Error(err)
			}
			logger.Trace(fmt.Sprintf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)))
		} else {
			logger.Error(fmt.Sprintf("Consumer error: %v (%v)\n", err, msg))
		}
	}
}

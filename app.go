/****
 * Copyright (c) 2020 Russia9
 */

package main

import (
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"kod-guildbot/bot"
	"kod-guildbot/lib"
	"os"
)

var logger = logrus.New()

func main() {
	logger.Out = os.Stdout

	logger.Info("Initializing CWBR-Bot")

	switch os.Getenv("CWBR_LOGLEVEL") {
	case "TRACE":
		logger.SetLevel(logrus.TraceLevel)
		break
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
		break
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
		break
	case "WARN":
		logger.SetLevel(logrus.WarnLevel)
		break
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
		break
	case "FATAL":
		logger.SetLevel(logrus.FatalLevel)
		break
	case "PANIC":
		logger.SetLevel(logrus.PanicLevel)
		break
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	dbHost := lib.GetEnvOrDefault("CWBR_DB_HOST", "localhost")
	dbPort := lib.GetEnvOrDefault("CWBR_DB_PORT", "3306")
	dbUser := lib.GetEnvOrDefault("CWBR_DB_USER", "cwbroker")
	dbPassword := lib.GetEnvOrDefault("CWBR_DB_PASS", "")
	dbName := lib.GetEnvOrDefault("CWBR_DB_NAME", "cwbroker")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		logger.Panic(db)
		return
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": lib.GetEnvOrDefault("CWBR_KAFKA_ADDRESS", "localhost"),
		"group.id":          "cw3",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		logger.Panic(err)
		return
	}

	err = bot.InitBot(os.Getenv("CWBR_BOT_TOKEN"), logger, consumer, db)
	if err != nil {
		logger.Panic(err)
	}
}

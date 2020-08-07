package lib

import (
	"database/sql"
	"gopkg.in/tucnak/telebot.v2"
	"os"
	"strconv"
)

func CheckPermsAndReply(needPerms int, bot *telebot.Bot, message *telebot.Message, db *sql.DB) bool {
	if !CheckPerms(message.Sender.ID, needPerms, db) {
		bot.Send(message.Chat, "Не хватает прав. Если это ошибка, пиши: @RussiaNine")
		return false
	}
	return true
}

func CheckPerms(id int, needPerms int, db *sql.DB) bool {
	if GetPerms(id, db) >= needPerms {
		return true
	}
	return false
}

func GetPerms(id int, db *sql.DB) int {
	var perms int
	permsRow := db.QueryRow("SELECT perms FROM users WHERE `id`=" + strconv.Itoa(id))
	err := permsRow.Scan(&perms)
	if err == nil {
		return perms
	} else {
		return -1
	}
}

func GetEnvOrDefault(name string, def string) string {
	if os.Getenv(name) == "" {
		return def
	}
	return os.Getenv(name)
}

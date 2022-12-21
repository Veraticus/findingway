package murult

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func EmojiFromStr(str string, emojis []*discordgo.Emoji) string {
	for _, e := range emojis {
		if e.Name == str {
			return fmt.Sprintf("<:%s>", e.APIName())
		}
	}
	return fmt.Sprintf(":question: (%s)", str)
}

func JobEmojiFromStr(str string, emojis []*discordgo.Emoji) string {
	for _, e := range emojis {
		if e.Name == str {
			return fmt.Sprintf("<:%s>", e.APIName())
		}
	}
	return fmt.Sprintf(":clown: (%s)", str)
}

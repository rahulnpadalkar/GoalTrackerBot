package customKeyboard

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func AddCustomKeys() tb.ReplyButton {

	allDone := tb.ReplyButton{Text: "All Done ✅"}
	return allDone

}

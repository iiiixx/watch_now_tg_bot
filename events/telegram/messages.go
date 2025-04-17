package telegram

const msgHelp = `I can save films for you and offer to watch later📀

+ to save a film, just send me the link🔗
+ to get a random film from your list, use the command /rnd🎰
(film will be removed after that🗑)`

const msgHello = "Yo, yo, yo! 🤙 \n\n" + msgHelp

const (
	msgUnknownCommand = "I'm confused💀"
	msgNoSavedPages   = "oops, nothing to show! you can send me a link🔗"
	msgSaved          = "saved it!🤝"
	msgAlreadyExists  = "this film is already on your list🐒"
)

package notifybot

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type Pusher interface {
	Do(interface{})
	SetTemplateBuildFunc(func() []string, func(name string) ([]byte, error)) error
	SetLogger(logger Logger)
}

type pusher struct {
	bot    *tg.BotAPI
	chatID int64
	temple Templater
	logger Logger
}

func InitPusher(chatId int64, token string) (Pusher, error) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	p := &pusher{
		bot:    bot,
		chatID: chatId,
	}
	return p, nil
}

func (p *pusher) Do(data interface{}) {

	text, err := p.temple.Do(data)
	if err != nil {
		p.logger.Error("Template run with error %s", err.Error())
		return
	}
	err = p.push(text)
	if err != nil {
		p.logger.Error("Send message with error %s", err.Error())
	}

}

func (p *pusher) SetTemplateBuildFunc(names func() []string, assets func(name string) ([]byte, error)) error {
	var err error
	p.temple, err = buildTemplates(assets, names)
	if err != nil {
		return err
	}
	return err
}

func (p *pusher) SetLogger(logger Logger) {
	p.logger = logger
}

func (p *pusher) push(txt string) error {
	var chattel tg.Chattable

	chattel = tg.NewMessage(p.chatID, strings.TrimRight(txt, "\n"))
	_, err := p.bot.Send(chattel)
	if err != nil {
		return err
	}

	return nil
}

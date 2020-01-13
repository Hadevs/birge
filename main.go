package main

import (
	"os"
	"log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var (
		port      = os.Getenv("PORT")       // sets automatically
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

  // here are buttons defined

  backBtn := tb.InlineButton{
    Unique: "back",
    Text:   "↩️ Назад"}

  enterBtn := tb.InlineButton{
    Unique: "enter",
    Text:   "🔑 Войти на биржу"}

  qualifyBtn := tb.InlineButton{
    Unique: "qualify",
    Text:   "🧧 Подать заявку"}

  infoBtn := tb.InlineButton{
    Unique: "info",
    Text:   "📃 Информация о бирже"}

  howToEnterBtn := tb.InlineButton{
    Unique: "howToEnter",
    Text:   "🗝 Как попасть на биржу?"}

  fuckedUpBtn := tb.InlineButton{
    Unique: "fuckedUp",
    Text:   "📆 Что будет, если я не уложусь в срок?"}

  whatProjectsBtn := tb.InlineButton{
    Unique: "whatProjects",
    Text:   "📑 Какие проекты предоставляет биржа?"}

  // currentProjectBtn := tb.InlineButton{
  //   Unique: "currentProject",
  //   Text:   "🛎 Мой текущий проект"}

  // showOffersBtn := tb.InlineButton{
  //   Unique: "showOffers",
  //   Text:   "📜 Посмотреть текущие предложения"}

  // askAdminBtn := tb.InlineButton{
  //   Unique: "askAdmin",
  //   Text:   "💡 Вопрос администрации"}

  // techSuppBtn := tb.InlineButton{
  //   Unique: "techSupp",
  //   Text:   "📦 Получить техническую помощь"}

  // redeemMilestoneProjectBtn := tb.InlineButton{
  //   Unique: "redeemMilestoneProject",
  //   Text:   "✅ Закрыть этап/проект"}

  // cancelProjectBtn := tb.InlineButton{
  //   Unique: "cancelProject",
  //   Text:   "❌ Отказаться от проекта"}

  // takeProjectBtn := tb.InlineButton{
  //   Unique: "takeProject",
  //   Text:   "❇️ Принять проект #1"}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

  b.Handle("/start", func(m *tb.Message) {
    inlineKeys := [][]tb.InlineButton{
      []tb.InlineButton{enterBtn, qualifyBtn},
      []tb.InlineButton{infoBtn}}

    b.Send(
      m.Sender,
      "Добро пожаловать в Swift Exchange! Пожалуйста, выберите следующий шаг:",
      &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

  b.Handle(&infoBtn, func(c *tb.Callback) {
    b.Respond(c, &tb.CallbackResponse{
        ShowAlert: false,
    })

    inlineKeys := [][]tb.InlineButton{
      []tb.InlineButton{howToEnterBtn},
      []tb.InlineButton{fuckedUpBtn,whatProjectsBtn},
      []tb.InlineButton{backBtn}}

    b.Send(c.Sender, `📃 Информация о бирже:

Swift Exchange - приватная биржа для доверенных разработчиков.

Мы берем на себя:

📩 Полное общение с заказчиками
💌 Профессиональную и техническую помощь в любом вопросе
📅 Поднятие вашего рейтинга, как разработчика, развитие личного бренда
📈 Постоянный поток проектов

Биржа забирает 5% с каждого проекта и выплачивает разработчику заработанные деньги сразу после принятия работ заказчиком. Способ оплаты обсуждается с каждым разработчиком отдельно.`,
    &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!")
	})

	b.Start()
}

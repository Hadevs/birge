package main

import (
	"os"
	"log"
  "fmt"
	tb "gopkg.in/tucnak/telebot.v2"
  "github.com/gomodule/redigo/redis"
)

func main() {
	var (
		port      = os.Getenv("PORT")       // sets automatically
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
    redisURL  = os.Getenv("REDIS_URL")
	)

  client, err := redis.DialURL(redisURL)
  if err != nil {
    log.Fatal(err)
  }
  defer client.Close()

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

  b.Handle("/whoami", func(m *tb.Message) {
    client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "whoami")
    client.Flush()
    client.Receive()
    b.Send(m.Sender, fmt.Sprintf("%d", m.Sender.ID))
  })

  b.Handle("/start", func(m *tb.Message) {
    client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "start")
    client.Flush()
    client.Receive()
    inlineKeys := [][]tb.InlineButton{
      []tb.InlineButton{enterBtn, qualifyBtn},
      []tb.InlineButton{infoBtn}}

    b.Send(
      m.Sender,
      "Добро пожаловать в Swift Exchange! Пожалуйста, выберите следующий шаг:",
      &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

  b.Handle(&infoBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "info")
    client.Flush()
    client.Receive()
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

  b.Handle(&howToEnterBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "howToEnter")
    client.Flush()
    client.Receive()

    inlineKeys := [][]tb.InlineButton{[]tb.InlineButton{backBtn}}

    b.Send(
      c.Sender,
      `🗝 Как попасть на биржу?:

Любой разработчик может попасть на биржу. Для этого достаточно нажать на соответствующую кнопку в начале диалога, рассказать пару слов о себе и с Вами свяжется администрация биржи для небольшого текстового интервью. Если вы уже работали на фрилансе, выполняли проекты и вы добросовестный разработчик, Вы обязательно попадете на биржу. Так же Вам нужно будет оплатить символический вступительный взнос, дабы убедиться в Ваших намерениях в размере 350р.`,
      &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

  b.Handle(&whatProjectsBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "whatProjects")
    client.Flush()
    client.Receive()

    inlineKeys := [][]tb.InlineButton{[]tb.InlineButton{backBtn}}

    b.Send(
      c.Sender,
      `📑 Какие проекты предоставляет биржа?:

Биржа предлагает любые проекты, связанные с языком Swift. В основном это iOS/macOS приложения. Так же будут задействованы и другие платформы.

Так же проект не будет выставлен на общее обозрение. Мы обращаемся к каждому разработчику по его внутреннему рейтингу, начиная с высокого. Если разработчику подходит проект - мы его передаем. Если нет - то обсуждаем данный проект уже со следующим разработчиком, пока проект не найдет своего исполнителя.`,
      &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

  b.Handle(&fuckedUpBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "fuckedUp")
    client.Flush()
    client.Receive()

    inlineKeys := [][]tb.InlineButton{[]tb.InlineButton{backBtn}}

    b.Send(
      c.Sender,
      `📆 Что будет, если я не уложусь в срок?:

Если Вы понимаете, что опаздываете со сдачей проекта на пару дней - ничего страшного, за это не будет никакого штрафа. Если же Вы понимаете, что опаздываете на дней 5 или больше, вы должны предупредить администрацию об этом за 1 неделю до сдачи проекта. Если состояние кода удовлетворительное, то Вам так же ничего не грозит. Просто снятие с проекта и назначение нового.

В иных ситуациях, администрация будет вынуждена попросить Вас выплатить штраф или заблокировать на бирже.`,
      &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
  })

  b.Handle(&backBtn, func(c *tb.Callback) {
    client.Send("GET", fmt.Sprintf("%s", c.Sender.ID))
    client.Flush()
    v, err := client.Receive()
    if err != nil {
      log.Print(err)
    }
    log.Print(v)
    switch v {
      case "info":
        client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "start")
        client.Flush()
        client.Receive()
        inlineKeys := [][]tb.InlineButton{
          []tb.InlineButton{enterBtn, qualifyBtn},
          []tb.InlineButton{infoBtn}}
        b.Send(
          c.Sender,
          "Добро пожаловать в Swift Exchange! Пожалуйста, выберите следующий шаг:",
          &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
        return
      case "fuckedUp", "whatProjects", "howToEnter":
        client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "info")
        client.Flush()
        client.Receive()
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

    Биржа забирает 5% от суммы с каждого проекта и выплачивает разработчику заработанные деньги сразу после принятия работ заказчиком. Способ оплаты обсуждается с каждым разработчиком отдельно.`,
        &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
        return
      default:
        client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "start")
        client.Flush()
        client.Receive()
        inlineKeys := [][]tb.InlineButton{
          []tb.InlineButton{enterBtn, qualifyBtn},
          []tb.InlineButton{infoBtn}}
        b.Send(
          c.Sender,
          "Добро пожаловать в Swift Exchange! Пожалуйста, выберите следующий шаг:",
          &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
        return
    }
  })

	b.Start()
}

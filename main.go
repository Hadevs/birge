package main

import (
	"os"
	"log"
  "fmt"
  "strings"
  "strconv"

	tb "gopkg.in/tucnak/telebot.v2"

  // "database/sql"
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "github.com/gomodule/redigo/redis"
)

var schema = `
  CREATE TABLE IF NOT EXISTS SEworker(
    id SERIAL PRIMARY KEY,
    tid TEXT,
    approved BOOLEAN,
    cpid INT
  );

  CREATE TABLE IF NOT EXISTS SEproject(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    difficulty INT,
    price INT,
    paid INT,
    progress INT
  );
`

type SEworker struct {
  Id int `db:"id"`
  Tid string `db:"tid"`
  Approved bool `db:"approved"`
  Cpid int `db:"cpid"`
}

type SEproject struct {
  Id int `db:"id"`
  Name string `db:"name"`
  Description string `db:"description"`
  Difficulty int `db:"difficulty"`
  Price int `db:"price"`
  Paid int `db:"paid"`
  Progress int `db:"progress"`
}

func parsePsqlElements(url string) (string, string, string, string, string) {
  split := strings.Split(url, "@")
  unamepwdsplit := strings.Split(split[0], "//")
  unamepwd := strings.Split(unamepwdsplit[1], ":")
  uname := unamepwd[0]
  pwd := unamepwd[1]
  urlportdbname := strings.Split(split[1], ":")
  link := urlportdbname[0]
  portdbname := strings.Split(urlportdbname[1], "/")
  port := portdbname[0]
  dbname := portdbname[1]
  return uname, pwd, link, port, dbname
}

func main() {
	var (
		port      = os.Getenv("PORT")       // sets automatically
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
    redisURL  = os.Getenv("REDIS_URL")
    psqlURL   = os.Getenv("DATABASE_URL")
    dbuname, dbpwd, dblink, dbport, dbname = parsePsqlElements(psqlURL)
    psqlInfo  = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s" +
    " sslmode=disable", dblink, dbport, dbuname, dbpwd, dbname)
    // Cuz I'm too lazy to do it the right way
    projectname  = ""
    projectdesc  = ""
    projectdiff  = 0
    projectprice = 0
	)

  fmt.Println(psqlInfo)

  db, err := sqlx.Connect("postgres", psqlInfo)
  if err != nil {
    log.Panic(err)
  }
  defer db.Close()
  db.MustExec(schema)

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

  showOffersBtn := tb.InlineButton{
    Unique: "showOffers",
    Text:   "📜 Посмотреть текущие предложения"}

  askAdminBtn := tb.InlineButton{
    Unique: "askAdmin",
    Text:   "💡 Вопрос администрации"}

  techSuppBtn := tb.InlineButton{
    Unique: "techSupp",
    Text:   "📦 Получить техническую помощь"}

  redeemMilestoneProjectBtn := tb.InlineButton{
    Unique: "redeemMilestoneProject",
    Text:   "✅ Закрыть этап/проект"}

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

  b.Handle(&enterBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "enter")
    client.Flush()
    client.Receive()

    inlineKeys := [][]tb.InlineButton{[]tb.InlineButton{showOffersBtn}}

    user := SEworker{}
    err := db.Get(&user, "SELECT * FROM SEworker WHERE tid=$1", c.Sender.ID)
    if err != nil || user.Approved != true {
      b.Send(c.Sender, `Сначала надо пройти собеседование, для этого нажми на "🧧 Подать заявку"`)
      return
    }
    projects := []SEproject{}
    db.Select(&projects, "SELECT * FROM SEproject ORDER BY id DESC")
    b.Send(c.Sender, fmt.Sprintf(`🔑 Войти на биржу:

Вы вошли на Swift Exchange. У вас сейчас %d открытых предложений по проектам.`, len(projects)),
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

  b.Handle(&qualifyBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "qualify0")
    client.Flush()
    client.Receive()

    b.Send(
      c.Sender,
      `🧧 Подать заявку:

Мы очень рады, что Вы решили попробовать себя в нашей бирже. Пожалуйста, напишите кратко о себе, своем опыте в разработке и проектах, с которыми Вы сталкивались. После этого, в самые кратчайшие сроки с Вами свяжется администрация для интервью в текстовом виде. Будем ждать! 😉`)
  })

  b.Handle(&askAdminBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "askAdmin0")
    client.Flush()
    client.Receive()

    b.Send(
      c.Sender,
      `💡 Вопрос администрации:

Введите Ваш запрос и он будет направлен администрации. Пожалуйста, старайтесь детальнее описать Вашу проблему. Запросы в формате "У меня проблема, помогите." рассматриваться не будут.`)
  })

  b.Handle(&techSuppBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "techSupp0")
    client.Flush()
    client.Receive()

    b.Send(
      c.Sender,
      `📦 Получить техническую помощь:

Вы должны описать проблему в проекте, с которой столкнулись. Чем больше информации Вы предоставите, тем быстрее получите ответ на Ваш вопрос. Мы стараемся помочь Вам в самый краткий срок.

Формат обращения:

1) Название проекта
2) Описание проблемы
3) Приложенные части кода, залитые на pastebin.com или Github Gist, скриншоты или видео, на которых видно проблему`)
  })

  b.Handle(&redeemMilestoneProjectBtn, func(c *tb.Callback) {
    client.Send("SET", fmt.Sprintf("%s", c.Sender.ID), "redeemMilestoneProject0")
    client.Flush()
    client.Receive()

    b.Send(
      c.Sender,
      `✅ Закрыть этап/проект:

Вы собираетесь закрыть проект или этап. Пожалуйста, заполните форму для закрытия:

1) Название проекта
2) Номер этапа, если закрываете этап
3) hash-номер коммита, который можно запускать для тестирования`)
  })

  b.Handle(&backBtn, func(c *tb.Callback) {
    client.Send("GET", fmt.Sprintf("%s", c.Sender.ID))
    client.Flush()
    v, err := client.Receive()
    if err != nil {
      log.Print(err)
    }
    position := fmt.Sprintf("%s", v)
    switch position {
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

  b.Handle(&showOffersBtn, func(c *tb.Callback) {
    projects := []SEproject{}
    db.Select(&projects, "SELECT * FROM SEproject ORDER BY id DESC")
    for _, project := range projects {
      b.Send(c.Sender, fmt.Sprintf(`– %s
**Задача:** %s
**Сложность:** %d | **Стоимость:** %d
`, project.Name, project.Description, project.Difficulty, project.Price))
    }
  })

  b.Handle("/approve", func(m *tb.Message) {
    tx := db.MustBegin()
    tx.MustExec(`INSERT INTO SEworker(tid, approved, cpid) VALUES ($1, true, 0)`, m.Payload)
    tx.Commit()
    b.Send(m.Sender, "Успешно добавлен новый пидераст, деньги мне плати блять")
  })

  b.Handle("/project", func(m *tb.Message) {
    client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "project0")
    client.Flush()
    client.Receive()
    b.Send(m.Sender, "Ну че, хуила, новый проект нашел для плебеев? Ну заполняй блять, гандон. Деньги мне плати блять")
  })

  b.Handle(tb.OnText, func(m *tb.Message) {
    client.Send("GET", fmt.Sprintf("%s", m.Sender.ID))
    client.Flush()
    v, err := client.Receive()
    position := fmt.Sprintf("%s", v)
    if err != nil {
      log.Print(err)
      b.Send(m.Sender, "К сожалению, что-то пошло не так. Пожалуйста, попробуйте позже. Администрация уже работает над решением проблемы!")
      return
    }
    admin := tb.User{73346375,"","","","",false}
    switch position {
      case "qualify0":
        b.Send(&admin, fmt.Sprintf("%d", m.Sender.ID))
        b.Forward(&admin, m)
        b.Send(m.Sender, "Спасибо, администрация получила Вашу заявку и в самое ближайшее время свяжется с вами в Telegram! Вернитесь в меню с помощью /start")
        return
      case "askAdmin0":
        b.Forward(&admin, m)
        b.Send(m.Sender, "Спасибо, администрация получила Ваш вопрос и в самое ближайшее время свяжется с вами в Telegram! Вернитесь в меню с помощью /start")
        return
      case "techSupp0":
        b.Forward(&admin, m)
        b.Send(m.Sender, "Спасибо, администрация получила Ваш запрос тех. поддержки и в самое ближайшее время свяжется с вами в Telegram! Вернитесь в меню с помощью /start")
        return
      case "redeemMilestoneProject0":
        b.Forward(&admin, m)
        b.Send(m.Sender, "Спасибо, администрация получила Ваш запрос закрытие этапа/проекта и в самое ближайшее время свяжется с вами в Telegram! Вернитесь в меню с помощью /start")
        return
      case "project0":
        projectname = m.Text
        client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "project1")
        client.Flush()
        client.Receive()
        b.Send(m.Sender, "Теперь пиши блять описание для своего ссаного проекта, хуила. Деньги мне плати блять")
        return
      case "project1":
        projectdesc = m.Text
        client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "project2")
        client.Flush()
        client.Receive()
        b.Send(m.Sender, "Теперь пиши блять насколько ахуенно сложный проект ты там придумал (1-5). Деньги мне плати блять")
      case "project2":
        projectdiff, err = strconv.Atoi(m.Text)
        if err != nil {
          b.Send(m.Sender, "Ты ебанутый блять? Пиши цифры блять, ЦИФРЫ СУКА, ЗНАЕШЬ ТАМ 1,2,3,4,5,6,7,8,9,0? НЕТ? ДЕБИЛ БЛЯТЬ")
          return
        }
        client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "project3")
        client.Flush()
        client.Receive()
        b.Send(m.Sender, "Теперь пиши блять сколько грошей (рублей) ты заплатишь плебсу, который это говно делать будет. Деньги мне плати блять")
      case "project3":
        projectprice, err = strconv.Atoi(m.Text)
        if err != nil {
          b.Send(m.Sender, "Ты ебанутый блять? Пиши цифры блять, ЦИФРЫ СУКА, ЗНАЕШЬ ТАМ 1,2,3,4,5,6,7,8,9,0? НЕТ? ДЕБИЛ БЛЯТЬ")
          return
        }
        fmt.Println(projectname)
        fmt.Println(projectdesc)
        fmt.Println(projectdiff)
        fmt.Println(projectprice)
        tx := db.MustBegin()
        tx.MustExec(`INSERT INTO SEproject(name, description, difficulty, price, paid, progress) VALUES ($1, $2, $3, $4, 0, 0)`, projectname, projectdesc, projectdiff, projectprice)
        tx.Commit()
        b.Send(m.Sender, "Поздравляю, долбаеб, все готово, проект теперь в списке, иди ищи плебсов, чтобы этого говно делали. Деньги мне плати блять")
        client.Send("SET", fmt.Sprintf("%s", m.Sender.ID), "start")
        client.Flush()
        client.Receive()
        return
      default:
        b.Send(m.Sender, "Я не понимаю обычный текст, нажмите /start")
    }
	})

	b.Start()
}

package core

const (
	translator  = "translator"
	editor      = "editor"
	errorString = "Произошла ошибка, попробуйте еще раз"
	notFound    = "record not found"
	markdown    = "Markdown"
)

const (
	strAddArticle      = "Добавить статью"
	strAddArticleShort = "/add"
	strCancel          = "Забудь" // TODO: emoji
	strDone            = "Дальше" // TODO: emoji
	strAddPlus         = "+ "
	strAddMinus        = "− "
)

var (
	jobMessages = map[string]string{
		translator: "Выбери переводчиков",
		editor:     "Выбери редакторов",
	}

	jobReplies = map[string]string{
		translator: "***Переводчики***",
		editor:     "***Редакторы***",
	}
)

const (
	replyDefault     = "Чего изволите, %s %s?"
	replyCancel      = "Окей! Что-нибудь еще?"
	replyNoUser      = "Не могу найти этого пользователя! Введи имя или ID ВК"
	replyNoChosen    = "Но ты ничего не выбрал!"
	replyCategoryJob = "Выбери категории"
	replyNoCategory  = "Не могу найти эту категорию! Введи верную"
	replyCategories  = "***Категории***"
)

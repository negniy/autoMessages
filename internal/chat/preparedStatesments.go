package chat

// Список приветствий
var hi = []string{
	"Привет!",
	"Здравствуйте!",
	"Добрый день!",
	"Доброе утро!",
	"Добрый вечер!",
	"Приветствую!",
	"Хей!",
	"Здарова!",
	"Приветики!",
	"Салют!",
	"Мои почтения!",
	"Как дела?",
	"Йо!",
	"Как жизнь?",
	"Здорово!",
	"Добро пожаловать!",
	"Рад тебя видеть!",
	"Хеллоу!",
	"Хай!",
	"Доброго времени суток!",
}

// Список прощаний
var bye = []string{
	"Пока!",
	"До свидания!",
	"Увидимся!",
	"До скорого!",
	"Прощай!",
	"Береги себя!",
	"Всего доброго!",
	"До встречи!", "Счастливо!", "До новых встреч!", "Удачи!", "Прощай!", "До завтра!", "Пока-пока!", "До скорой встречи!", "Всего хорошего!", "Хорошего дня!", "До связи!", "Будь здоров!", "Увидимся позже!",
}

// Список предзаготовленных фраз
var phrases = []string{
	"Это была бы смешная шутка, если бы не было так грустно.",
	"Кто-то сегодня встал не с той ноги?",
	"Если бы у меня был выбор, я бы выбрал пиццу!",
	"Почему программисты любят тёмную тему? Чтобы скрыть свои ошибки!",
	"Это звучит как начало анекдота.",
	"Сегодня так солнечно, настроение тоже на высоте!",
	"Кажется, скоро пойдёт дождь.",
	"На улице ужасно холодно, нужно укутаться!",
	"Погода как-то не радует в последнее время.",
	"Слышал, завтра будет снег?",
	"Сколько ещё учиться? Кажется, бесконечно.",
	"Экзамен завтра, а я тут сижу и размышляю.",
	"Как успеть всё выучить за одну ночь?",
	"Учёба учёбой, но отдыхать тоже нужно!",
	"Кажется, лекции стали интереснее, или это просто я?",
	"Сегодня на работе всё шло наперекосяк.",
	"Кажется, дедлайн уже наступил, а я только начал.",
	"Работа дома — это отдых или всё-таки работа?",
	"Кажется, нужно обсудить это с командой.",
	"Как сделать рабочий день более продуктивным?",
	"Хочу заняться чем-то новым, что посоветуешь?",
	"Может, попробовать научиться рисовать?",
	"Сегодня утром пробежал 5 километров, чувствую себя отлично!",
	"Я решил освоить новую игру, это будет весело!",
	"Может, пора вернуться к старому хобби?",
	"Не могу поверить, что это сработало!",
	"Кажется, мы близки к решению!",
	"Это точно поможет?",
	"Давай сделаем что-то необычное.",
	"Мне нравится, как это звучит.",
	"А как ты это представляешь?",
	"Это вызывает улыбку!",
	"Я думаю, это станет популярным.",
	"Как это объяснить попроще?",
	"Сложно, но возможно!",
	"Ты точно уверен?",
	"Это удивительно!",
	"Это была неожиданная идея.",
	"Звучит странно, но возможно.",
	"Почему бы не попробовать?",
	"Это требует доработки.",
	"Кажется, это полезно!",
	"Можем ли мы быть уверены?",
	"Как это можно проверить?",
	"Это просто совпадение?",
	"Надо подумать ещё раз.",
	"С этим нужно разобраться.",
	"Мне кажется, это важно.",
	"Что думаешь по этому поводу?",
	"Это стоит внимания.",
	"Давай уточним детали.",
}

// Список смайликов
var emojis = []string{
	"😊", "😉", "😅", "😂", "🤔", "🙃", "😐", "😇", "😎", "😭",
	"👍", "👏", "🙌", "💡", "❓", "❗", "🤷", "🤷‍♂️", "🤷‍♀️", "🙄",
}

// Список дополнительных слов
var inserts = []string{
	"Эмм...", "Пиздец", "Ура", "Капец", "Бля", "Юху", "Треш", "Ес", "Кул", "Cool",
}

//Список путей к гс
var voices = []string{
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\voices\\1.ogg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\voices\\2.ogg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\voices\\3.ogg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\voices\\4.ogg",
}

//Список путей к фото
var pics = []string{
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\pics\\1.jpg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\pics\\2.jpg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\pics\\3.jpg",
	"C:\\Users\\dstep\\OneDrive\\Рабочий стол\\Go\\autoMessages\\pics\\4.jpg",
}
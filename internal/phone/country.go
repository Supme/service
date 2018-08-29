package phone

var country = map[string]string{
	`7`:    `Россия | Казахстан`,
	`61`:   `Австралия`,
	`43`:   `Австрия`,
	`994`:  `Азербайджан`,
	`355`:  `Албания`,
	`213`:  `Алжир`,
	`1907`: `Аляска`,
	`1264`: `Ангилья`,
	`244`:  `Ангола`,
	`376`:  `Андорра`,
	`1268`: `Антигуа и Барбуда`,
	`599`:  `Антильские острова`,
	`853`:  `Аомынь (Макао)`,
	`54`:   `Аргентина`,
	`374`:  `Армения`,
	`297`:  `Аруба`,
	`93`:   `Афганистан`,
	`1242`: `Багамские острова`,
	`880`:  `Бангладеш`,
	`1246`: `Барбадос`,
	`973`:  `Бахрейн`,
	`375`:  `Беларусь`,
	`501`:  `Белиз`,
	`32`:   `Бельгия`,
	`229`:  `Бенин`,
	`1441`: `Бермудские острова`,
	`95`:   `Бирма (Мьянма)`,
	`359`:  `Болгария`,
	`591`:  `Боливия`,
	`387`:  `Босния и Герцеговина`,
	`267`:  `Ботсвана`,
	`55`:   `Бразилия`,
	`673`:  `Бруней`,
	`226`:  `Буркина-Фасо`,
	`257`:  `Бурунди`,
	`975`:  `Бутан`,
	`678`:  `Вануату`,
	`44`:   `Великобритания`,
	`36`:   `Венгрия`,
	`58`:   `Венесуэла`,
	`1340`: `Виргинские острова (Американские)`,
	`1284`: `Виргинские острова (Британские)`,
	`84`:   `Вьетнам`,
	`241`:  `Габон`,
	`1808`: `Гавайские острова`,
	`509`:  `Гаити`,
	`592`:  `Гайана`,
	`220`:  `Гамбия`,
	`233`:  `Гана`,
	`590`:  `Гваделупа`,
	`502`:  `Гватемала`,
	`594`:  `Гвиана Французская`,
	`224`:  `Гвинея`,
	`245`:  `Гвинея-Бисау`,
	`49`:   `Германия`,
	`350`:  `Гибралтар`,
	`852`:  `Гонгконг`,
	`504`:  `Гондурас`,
	`1473`: `Гренада`,
	`299`:  `Гренландия`,
	`30`:   `Греция`,
	`995`:  `Грузия`,
	`1671`: `Гуам`,
	`45`:   `Дания`,
	`243`:  `Демократическая Республика Конго (Заир)`,
	`253`:  `Джибути`,
	`1767`: `Доминика`,
	`1809`: `Доминиканская Республика`,
	`20`:   `Египет`,
	`260`:  `Замбия`,
	`263`:  `Зимбабве`,
	`972`:  `Израиль`,
	`91`:   `Индия`,
	`62`:   `Индонезия`,
	`962`:  `Иордания`,
	`964`:  `Ирак`,
	`98`:   `Иран | Коморские острова`,
	`353`:  `Ирландия`,
	`354`:  `Исландия`,
	`34`:   `Испания`,
	`39`:   `Италия`,
	`967`:  `Йемен`,
	`238`:  `Кабо Верде`,
	`1345`: `Каймановы острова`,
	`855`:  `Камбоджа`,
	`237`:  `Камерун`,
	`3428`: `Канарские острова`,
	`974`:  `Катар`,
	`254`:  `Кения`,
	`357`:  `Кипр`,
	`686`:  `Кирибати`,
	`86`:   `Китай`,
	`850`:  `КНДР (Северная Корея)`,
	`57`:   `Колумбия`,
	`269`:  `Коморские острова`,
	`242`:  `Конго`,
	`3395`: `Корсика`,
	`506`:  `Коста-Рика`,
	`225`:  `Кот-д"Ивуар`,
	`53`:   `Куба`,
	`965`:  `Кувейт`,
	`996`:  `Кыргызстан`,
	`856`:  `Лаос`,
	`371`:  `Латвия`,
	`266`:  `Лесото`,
	`231`:  `Либерия`,
	`961`:  `Ливан`,
	`218`:  `Ливия`,
	`370`:  `Литва`,
	`423`:  `Лихтенштейн`,
	`352`:  `Люксембург`,
	`230`:  `Маврикий`,
	`222`:  `Мавритания`,
	`261`:  `Мадагаскар`,
	`389`:  `Македония`,
	`265`:  `Малави`,
	`60`:   `Малайзия`,
	`223`:  `Мали`,
	`960`:  `Мальдивские острова`,
	`356`:  `Мальта`,
	`212`:  `Марокко`,
	`596`:  `Мартиника`,
	`52`:   `Мексика`,
	`691`:  `Микронезия`,
	`258`:  `Мозамбик`,
	`373`:  `Молдова`,
	`377`:  `Монако`,
	`976`:  `Монголия`,
	`1664`: `Монтсеррат`,
	`264`:  `Намибия`,
	`674`:  `Науру`,
	`977`:  `Непал`,
	`227`:  `Нигер`,
	`234`:  `Нигерия`,
	`31`:   `Нидерланды`,
	`505`:  `Никарагуа`,
	`64`:   `Новая Зеландия`,
	`687`:  `Новая Каледония`,
	`47`:   `Норвегия`,
	`672`:  `Норфолк остров`,
	`971`:  `ОАЭ (Объединенные Арабские Эмираты)`,
	`968`:  `Оман`,
	`92`:   `Пакистан`,
	`6809`: `Палау`,
	`507`:  `Панама`,
	`675`:  `Папуа-Новая Гвинея`,
	`595`:  `Парагвай`,
	`51`:   `Перу`,
	`48`:   `Польша`,
	`351`:  `Португалия`,
	`1787`: `Пуэрто-Рико`,
	`262`:  `Реюньон`,
	`250`:  `Руанда`,
	`40`:   `Румыния`,
	`503`:  `Сальвадор`,
	`685`:  `Самоа (Западное Самоа)`,
	`378`:  `Сан-Марино`,
	`239`:  `Сан-Томе и Принсипи`,
	`966`:  `Саудовская Аравия`,
	`268`:  `Свазиленд`,
	`1670`: `Северные Марианские острова`,
	`248`:  `Сейшельские острова`,
	`221`:  `Сенегал`,
	`1784`: `Сент-Винсент и Гренадины`,
	`1869`: `Сент-Китс и Невис`,
	`1758`: `Сент-Люсия`,
	`381`:  `Сербия и Черногория`,
	`65`:   `Сингапур`,
	`963`:  `Сирия`,
	`421`:  `Словакия`,
	`386`:  `Словения`,
	`252`:  `Сомали`,
	`249`:  `Судан`,
	`597`:  `Суринам`,
	`1`:    `США (Соединенные Штаты Америки) | Канада`,
	`232`:  `Сьерра-Леоне`,
	`992`:  `Таджикистан`,
	`66`:   `Таиланд`,
	`886`:  `Тайвань`,
	`255`:  `Танзания`,
	`1649`: `Теркс и Кайкос острова`,
	`228`:  `Того`,
	`690`:  `Токелау`,
	`676`:  `Тонга`,
	`1868`: `Тринидад и Тобаго`,
	`216`:  `Тунис`,
	`993`:  `Туркменистан`,
	`90`:   `Турция`,
	`256`:  `Уганда`,
	`998`:  `Узбекистан`,
	`380`:  `Украина`,
	`598`:  `Уругвай`,
	`298`:  `Фарерские острова`,
	`679`:  `Фиджи`,
	`63`:   `Филиппины`,
	`358`:  `Финляндия`,
	`33`:   `Франция`,
	`689`:  `Французская Полинезия`,
	`385`:  `Хорватия`,
	`236`:  `ЦАР (Центрально-Африканская Республика)`,
	`235`:  `Чад`,
	`420`:  `Чехия`,
	`56`:   `Чили`,
	`41`:   `Швейцария`,
	`46`:   `Швеция`,
	`94`:   `Шри-Ланка`,
	`593`:  `Эквадор`,
	`240`:  `Экваториальная Гвинея`,
	`291`:  `Эритрея`,
	`372`:  `Эстония`,
	`251`:  `Эфиопия`,
	`27`:   `ЮАР`,
	`82`:   `Южная Корея`,
	`1876`: `Ямайка`,
	`81`:   `Япония`,
}

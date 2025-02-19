# API онлайн библиотека песен
Данный сервис реализован на языке Go с использованием библиотеки Gin Gonic. Для работы с PostgreSQL использовался драйвер pgx, для миграций lib/pq. Для миграции использовалась библиотека goose, миграции выполняются автоматически при запуске сервиса. Тип операции зависит от определенного запроса. Для данных из файла конфигурации .env используется GoDotEnv.
Данные принимаются в JSON формате, далее сохраняются в базу данных PostgreSQL. Вывод данных происходит тоже в формате JSON. Реализованы все основные операции сохранения, изменения, получения и удаления.
Операции выполняются вместе со библиотекой песен. При сохранении песни происходит запрос в API.
Реализована спецификация на API в формате Swagger 2.0 с подходом code-first. Документация Swagger после запуска будет доступна по ссылке: http://localhost:8000/swagger/index.html#/
В Makefile прописаны возможные варианты запуска API и миграции.
## Есть 2 способа запуска микросервиса:
### 1. Локально.
   Необходимо в файле конфигурации .env указать
   ```
   HOST=localhost
   PORT=5436
   ```
   Запустить контейнер postgres
   ```
   docker run --name=db -e POSTGRES_PASSWORD='54321' -p 5436:5432 -d postgres
   ```
   Ввести в консоль команду
   ```
   go run cmd/main.go
   ```
   При необходимости, можно заменить названия контейнера базы данных, пароль и порты. Соответствующие параметры для такого же изменения находятся в .env. 
   Также, для отката миграции можно выполнить команду, для этого нужно установить goose (информация по ссылке: https://github.com/pressly/goose)
   ```
   make migrate-down
   ```
### 2. Локально при использовании docker-compose.
   
   Необходимо в файле конфигурации .env указать (По умолчанию в проекте стоит такое значение)
   ```
   HOST=db
   PORT=5432
   ```
   Ввести в консоль команду
   ```
   make build && make run
   ```
## Пользование сервисом
### Для получения списка данных библиотеки с фильтрацией по всем полям и пагинацией
```
curl --location --request GET 'http://localhost:8000/api/songs?sort={artist}&order={asc}&page={1}&name={muse}' \
--header 'Content-Type: application/json' \
--data ''
```
Данные в запросе указаны для примере. Вместо sort вводится желаемая сортировка списка. Если поле оставить пустым, по умолчанию будет применена сортировка по исполнителю.
Вместо order нужно ввести желаемое  направление списка. Если поле оставить пустым, по умолчанию будет применена сортировка по исполнителю. 
Вместо page нужно ввести желаемое количество страниц с данными. Установлено отображение 10 записей на 1 страницу. Если поле оставить пустым, по умолчанию будет значение 1. 
Далее можно вводить name, group, text, releasedate, link с указанием данных для фильтрации по всем полям. 
После успешного запроса будет выведен список сохранненых в базу данных песен.

#### Для добавления песни необходимо ввести запрос
```
curl --location  --request POST 'http://localhost:8000/api/songs' \
--header 'Content-Type: application/json' \
--data '           {
                "group": "Muse",
                "name": "Supermassive Black Hole"
            }'
```
В полях для добавления нужно ввести исполнителя и название песни. По этим данным будет выполняться описанный в ТЗ запрос и получение данных о дате релиза, тексте песни и ссылке на песню.
После успешного запроса будет выведен идентификатор сохранненой песни. 

#### Для обновления песни необходимо ввести запрос
```
curl --location --request PUT 'http://localhost:8000/api/songs/{id}' \
--header 'Content-Type: application/json' \
--data '           {
                "group": "Muse",
                "name": "Supermassive Black Hole",
                "date": "16.07.2006",
                "text": "Ooh baby, don'\''t you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
                "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
            }'
```
Для исправления можно ввести любое нужное количество полей в указанном порядке. В поле запроса вместо id указать идентификатор фильма для изменения.

#### Для удаления песни необходимо ввести запрос
```
curl --location --request DELETE 'http://localhost:8000/api/songs/id' \
--header 'Content-Type: application/json' \
--data ''
```
В поле запроса вместо id указать идентификатор фильма для удаления.

#### Для получения текста песни с пагинацией по куплетам
```
curl --location --request GET 'http://localhost:8000/api/songs/get-text?name={muse}&begin={1}&end={1}' \
--header 'Content-Type: application/json' \
--data ''
```
В поле запроса вместо name указать название песни для отображения текста песни. Вместо begin и end желаемый диапазон куплетов, в примере указан 1 куплет.
## Обработка ошибок
Для различных методов и вызовов функций реализована обработка ошибок, в зависимости от категории ошибки, выдается текст и код ошибки. Присуствуют коды 4хх и 5хх.

# Проект: сокращатель ссылок

В проекте присутствует возможность выбора хранилища. При запуске проекта необходимо указать параметр.

Запуск проекта с сохранением в памяти приложения:
```
make storage=cache
```

Запуск проекта с сохранением в базе данных:
```
make storage=postgres
```
В случае запуска с сохранением в БД, запускаются так-же контейнеры с БД и с миграциями.

Запуск тестов:
```
make test
```

Пересобрать проект:
```
make rebuild storage=*вариант хранилища*
```

Логи приложения записываются в файл:
```logs/server.log```

Параметры для запуска сервера и БД указываются в файле ".env". Пример .env:
```
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DATABASE_URL=postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

LISTEN_TYPE=port
BIND_IP=0.0.0.0
PORT=8080
```

Документация к проекту:
```
http://localhost:8080/swagger/index.html
```

# Simple gRPC service

## ТЗ
1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей
2. Реализовать gRPC сервис на основе proto файла на Go
3. Для хранения данных использовать PostgreSQL
4. на запрос получения списка пользователей данные будут кешироваться в redis на минуту и брать из редиса
5. При добавлении пользователя делать лог в clickHouse
6. Добавление логов в clickHouse делать через очередь Kafka

## Перед запуском
Установить:
- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://go.dev/doc/install)
- [gPRC](https://grpc.io/docs/languages/go/quickstart/)
- gcc ([для Windows](https://jmeubank.github.io/tdm-gcc/download/))

## Быстрый запуск

Команда запускает в докерах все сервисы и выводит результат в консоль
```bash
docker-compose -f .\docker-compose-quick.yml up --build
```

## Запуск

Запускаем базу

```
docker-compose up
```

Запускаем сервер.

```
cd service/server
go run server/server.go
```

Запускаем сервис для логирования.

```
cd service/logger
go run logger/logger.go
```

Отдельно запускаем клиент с тестовыми запросами

```
cd service/server
go run client/client.go
```

## Разработка


После изменения `.proto`-файла, генерируем код командой:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user/user.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/logger/logger.proto
```
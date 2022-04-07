# Simple gRPC service

## ТЗ
1. Описать proto файл с сервисом из 3 методов: добавить пользователя, удалить пользователя, список пользователей
2. Реализовать gRPC сервис на основе proto файла на Go
3. Для хранения данных использовать PostgreSQL
4. на запрос получения списка пользователей данные будут кешироваться в redis на минуту и брать из редиса
5. При добавлении пользователя делать лог в clickHouse
6. Добавление логов в clickHouse делать через очередь Kafka

## Запуск

Генерируем код из `.proto`-файла.

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto
```

Запускаем сервер.

```
go run server/server.go
```

Запускаем клиент с тестовыми запросами


```
go run client/client.go
```
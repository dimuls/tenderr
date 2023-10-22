# tenderr

Состоит из двух сервисов:
* tenderr-classifier - классификатор логов
* tenderr-operator - рабочее место оператора

## Содержание репазитория

* ./cmd - точки входа в программы
* ./entity - сущности используемые в коде
* ./mock - мок-реализации (на данные момент MessageSender - отправка сообщений о разрешении проблемы)
* ./postges - реализация хранилища на базе СУБД Postgres
* ./service/classifier - ядро сервиса tenderr-classifier
* ./service/operator - ядро сервиса tenderr-operator

## Локальный запуск системы (в linux)

Поднять докеры с postgres, clickhouse и grafana 
```bash
docker-compose up -d
```

Сборка фронтенда tenderr-classifier
```bash
cd ./services/classifier/ui
npm install
npm run build
```

Сборка фронтенда tenderr-operator
```bash
cd ./services/operator/ui
npm install
npm run build
```

Запуск tenderr-classifer
```bash
go build -o ./bin/ ./cmd/tenderr-classifier && ./bin/tenderr-classifier -config ./services/classifier/config.example.yaml
```

Запуск tenderr-operator
```bash
go build -o ./bin/ ./cmd/tenderr-operator && ./bin/tenderr-operator -config ./services/operator/config.example.yaml
```

## Локальные адреса

* [grafana](http://localhost:3000)
* [tenderr-classifer](http://localhost:8080)
* [tenderr-operator](http://localhost:8081)

## Примеры запросов к tenderr-operator

Создание ошибки от пользователя
```bash
curl -X POST http://localhost:8081/api/user-errors -d '{"elementId":"1946d729e2ddc19eeb747ad19561f8f9","message":"Не работает кнопка продолжить","contact":{"type":"telegram","data":"@dimuls"}}'
```

Подписка на оповещение на решение проблемы:
```bash
curl -X POST http://localhost:8081/api/error-resolve-waiter -d '{"errorNotificationId":"a56892b8-eaf9-435d-af7e-28235553c013","contact":{"type":"telegram","data":"@dimuls"}}'
```
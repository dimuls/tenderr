# tenderr

Состоит из двух сервисов:
* tenderr-classifier - классификатор логов
* tenderr-operator - рабочее место оператора

## Локальный запуск системы

Поднять докеры с postgres, clickhouse и grafana 
```bash
docker-compose up -d
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
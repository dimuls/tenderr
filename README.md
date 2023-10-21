# tenderr

Состоит из двух сервисов:
* tenderr-classifier - классификатор логов
* tenderr-operator - рабочее место оператора

## Локальный запуск системы

Поднимите докеры с postgres, clickhouse и grafana 
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
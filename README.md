# wallet_controller

Небольшой сервис для работы с кошельком

## Описание
Позволяет:
- получить данные о кошельке по айди
- совершить операцию

Тестовый config.env для локального запуска (не в докере)
```azure
DB_HOST=localhost
DB_PORT=5432
DB_NAME=wallet_controller_db
DB_USERNAME=postgres
DB_PASSWORD=postgres

IP_ADDRESS=localhost
API_PORT=8080
```

### Запуск через Docker Compose

```bash
# Клонируйте репозиторий
git clone https://github.com/Ferginin/wallet_controller.git
cd wallet_controller

# Запустите сервисы
docker-compose up -d
```

Сервис будет доступен по адресу: **http://localhost:8080**

### Эндпоинты

```bash
http://localhost:8080/health
# Ожидаемый ответ:
# {"status":"ok"}

http://localhost:8080/api/v1/wallets/{UUID}
# Ожидаемый ответ:
# {"id": "33333333-3333-3333-3333-333333333333","balance": 2891000}

http://localhost:8080/api/v1/wallet
# примерное тело запроса:
#{
#    "wallet_id": "33333333-3333-3333-3333-333333333333",
#    "operation_type": "DEPOSIT",
#    "amount": 500
#}

# Ожидаемый ответ:
# {
#    "wallet": {
#        "id": "33333333-3333-3333-3333-333333333333",
#        "balance": 2941000
#    }
#}
```

## Тесты

для части тестов (wallet_repository_test.go) нужно создать бд wallet_test в postgresql
- владелец: postgres
- пароль: postgres
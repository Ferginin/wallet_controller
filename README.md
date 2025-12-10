# wallet_controller

Небольшой сервис для работы с кошельком

## Описание
Позволяет:
- получить данные о кошельке по айди
- совершить операцию

Тестовый config.env
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

### Проверка работы

```bash
# Health check
curl http://localhost:8080/health

# Ожидаемый ответ:
# {"status":"ok"}
```
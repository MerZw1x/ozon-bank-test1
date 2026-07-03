## Миграции базы данных

### Управление миграциями через Goose

Для управления схемой БД используется `goose`. Все миграции лежат в папке `migration/`.

### Установка goose

``` bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Доступные команды

```bash
# Посмотреть статус миграций
sh scripts/migrations.sh --status

# Создать новую миграцию
sh scripts/migrations.sh --new <название_миграции>
 
# Накатить все миграции
sh scripts/migrations.sh --up

# Откатить последнюю миграцию
sh scripts/migrations.sh --down

# Накатить до конкретной версии
sh scripts/migrations.sh --up <версия>

# Откатить до конкретной версии
sh scripts/migrations.sh --down <версия>

```

### Запуск миграция

``` bash
# Запускаем docker compose
docker compose up -d test_db

# Накатываем миграции
sh scripts/migrations.sh --up

# Проверяем статус
sh scripts/migrations.sh --status
```
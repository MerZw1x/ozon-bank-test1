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
sh migration/migrations.sh --status

# Создать новую миграцию
sh migration/migrations.sh --new <название_миграции>
 
# Накатить все миграции
sh migration/migrations.sh --up

# Откатить последнюю миграцию
sh migration/migrations.sh --down

# Накатить до конкретной версии
sh migration/migrations.sh --up <версия>

# Откатить до конкретной версии
sh migration/migrations.sh --down <версия>

```

### Запуск миграция

``` bash
# Запускаем docker compose
docker compose up -d db

# Накатываем миграции
sh migration/migrations.sh --up

# Проверяем статус
sh migration/migrations.sh --status
```
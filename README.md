# Файлы для итогового задания
# go_final_project

Этот проект представляет собой веб-сервер для управления задачами (TODO-list).

## Функционал

- Создание, редактирование и удаление задач
- Просмотр списка задач
- Отметка задач как выполненных
- Повторяющиеся задачи
- API для работы с задачами

## Сборка проекта

1. Убедитесь, что у вас установлен Go версии 1.x или выше
2. Клонируйте репозиторий
3. Перейдите в директорию проекта
4. Выполните команду: `go build -o scheduler`

## Запуск

1. Создайте файл `.env` в корневой директории проекта (опционально)
2. Запустите исполняемый файл: `./scheduler`
3. Сервер будет доступен по адресу: http://localhost:7540/

## Запуск тестов

Для запуска тестов выполните следующую команду в корневой директории проекта:




go test ./tests


Это запустит все тесты, находящиеся в директории `tests`.

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.
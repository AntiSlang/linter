# Go Log Linter

Линтер для проверки лог-сообщений

Версия Go: 1.25

1. Лог-сообщения должны начинаться со строчной буквы
2. Лог-сообщения должны быть только на английском языке
3. Лог-сообщения не должны содержать спецсимволы или эмодзи
4. Лог-сообщения не должны содержать потенциально чувствительные данные

Сборка осуществляется через Makefile. Возможно собрать плагин для golangci-lint, либо как отдельный бинарник

Установка плагина в проект:
1. Собрать плагин: `make build-plugin`
2. Создать файл `.golangci.yml` с содержимым:
```
version: "2"

linters:
  enable:
    - loglinter
  settings:
    custom:
      loglinter:
        type: goplugin
        path: [путь до linter.so]
        description: "Checks logging messages"
```
3. Запустить: `golangci-lint run --config .golangci.yml ./...`

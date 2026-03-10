# Go Log Linter

Линтер для проверки лог-сообщений

Версия Go: 1.25

Правила:
1. Лог-сообщения должны начинаться со строчной буквы
2. Лог-сообщения должны быть только на английском языке
3. Лог-сообщения не должны содержать спецсимволы или эмодзи
4. Лог-сообщения не должны содержать потенциально чувствительные данные

Сборка осуществляется через Makefile. Возможно собрать плагин для golangci-lint, либо как отдельный бинарник

Установка плагина в проект:
1. Собрать плагин: `make build-plugin`
2. Создать файл `.golangci.yml` с содержимым (настройки можно менять:
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
        settings:
          check-lowercase: true
          sensitive-words: "password,token,api key,configtest"
```

Поддерживаемые настройки:
- checkLowercase - включить / отключить 1 правило (по умолчанию включено)
- checkEnglish - включить / отключить 2 правило (по умолчанию включено)
- checkSpecials - включить / отключить 3 правило (по умолчанию включено)
- checkSensitive - включить / отключить 4 правило (по умолчанию включено)
- sensitiveWords - список слов для 4 правила (чувствительные данные)

3. Запустить: `golangci-lint run --config .golangci.yml ./...`

Пример:
![Скриншот](https://i.ibb.co/qLNfkdjB/example.png)
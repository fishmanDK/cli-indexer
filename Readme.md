# Запуск приложения из GitHub

Этот документ описывает шаги, необходимые для запуска приложения из репозитория GitHub.

## 1. Клонирование репозитория

- **Клонируйте репозиторий**: Используйте команду в терминале для клонирования репозитория на ваш компьютер:
```bash
  git clone https://github.com/fishmanDK/rest-entangle.git
```

## 2. Перейдите в каталог проекта
```bash
  cd ПАПКА_С_ПРОЕКТОМ
```

## 3. Запуск проекта
```bash
  go run main.go indexer run --rpc=<URL_BLOCKCHAIN_RPC> --start=<NUMBER_OF_BLOCK> --out=<PATH_TO_OUTPUT_FILE>
```
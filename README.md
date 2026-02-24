# 🐳 MockWhale

**MockWhale** — это легковесный и быстрый self-hosted сервис на Go для создания мок-эндпоинтов с удобным Dashboard.

## 🚀 Возможности
* **Dynamic Mocking**: Создание API-заглушек через встроенный дашборд или REST API.
* **Persistent Storage**: Все данные хранятся в локальной базе SQLite.
* **Pure Go**: Никаких внешних зависимостей (CGO-free).
* **Developer Friendly**: Встроенная админка по адресу `/_dashboard`.

## 🛠️ Технологический стек
* **Backend**: Golang
* **Database**: SQLite (ModernC)
* **Frontend**: Vanilla HTML/CSS (embedded in Go)

## 📦 Быстрый старт
1. Склонируйте репозиторий.
2. Запустите сервер:
   ```bash
   go run cmd/api/main.go
Откройте дашборд: http://localhost:3000/_dashboard

Разработано Chavs23 🇰🇿

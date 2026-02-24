# 🐳 MockWhale

![Go Version](https://img.shields.io/github/go-mod/go-version/Chavs23/mockwhale?style=flat-square&color=00ADD8)
![Database](https://img.shields.io/badge/database-SQLite-003B57?style=flat-square&logo=sqlite)
![Status](https://img.shields.io/badge/status-active-success?style=flat-square)

**MockWhale** — это легковесный и быстрый self-hosted сервис на Go для создания мок-эндпоинтов с удобным Dashboard.

## 🚀 Возможности
* **Dynamic Mocking**: Создание API-заглушек через встроенный дашборд или REST API.
* **Persistent Storage**: Все данные хранятся в локальной базе SQLite.
* **Pure Go**: Никаких внешних зависимостей (CGO-free).
* **Developer Friendly**: Встроенная админка по адресу `/_dashboard`.
* **JSON Editor**: Удобное поле для вставки ответов любой сложности.

## 🛠️ Технологический стек
* **Backend**: Golang
* **Database**: SQLite (ModernC)
* **Frontend**: Vanilla HTML/CSS (embedded in Go)

## 📦 Быстрый старт
1. **Склонируйте репозиторий:**
   ```bash
   git clone [https://github.com/Chavs23/mockwhale.git](https://github.com/Chavs23/mockwhale.git)
   Запустите сервер:
    Bash

    go run cmd/api/main.go

    Откройте дашборд: http://localhost:3000/_dashboard

Разработано Chavs23 🇰🇿

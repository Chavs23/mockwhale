package database

import (
"database/sql"
"os"
"fmt"
_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
db, err := sql.Open("sqlite", "./mockwhale.db")
if err != nil {
return nil, err
}

migration, err := os.ReadFile("migrations/001_init.sql")
if err != nil {
return nil, fmt.Errorf("не удалось прочитать файл миграции: %v", err)
}

_, err = db.Exec(string(migration))
if err != nil {
return nil, fmt.Errorf("ошибка при выполнении миграции: %v", err)
}

return db, nil
}

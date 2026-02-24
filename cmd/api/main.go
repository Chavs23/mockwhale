package main

import (
"database/sql"
"encoding/json"
"fmt"
"log"
"net/http"
"github.com/Chavs23/mockwhale/internal/database"
)

type MockRequest struct {
Path         string `json:"path"`
Method       string `json:"method"`
ResponseBody string `json:"response_body"`
StatusCode   int    `json:"status_code"`
}

func main() {
db, err := database.InitDB()
if err != nil {
log.Fatalf("Ошибка БД: %v", err)
}
defer db.Close()

seedMock(db)

port := ":3000"
fmt.Printf("🐳 MockWhale Dashboard: http://localhost%s/_dashboard\n", port)

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 1. Дашборд (Интерфейс)
if r.Method == http.MethodGet && r.URL.Path == "/_dashboard" {
renderDashboard(w, db)
return
}

// 2. Создание мока через Форму (для людей)
if r.Method == http.MethodPost && r.URL.Path == "/_create_web" {
path := r.FormValue("path")
method := r.FormValue("method")
status := r.FormValue("status")
// По дефолту для веба создаем пустой JSON объект
query := `INSERT INTO mock_endpoints (path, method, response_body, status_code) VALUES (?, ?, ?, ?)`
db.Exec(query, path, method, `{"status": "ok"}`, status)
http.Redirect(w, r, "/_dashboard", http.StatusSeeOther)
return
}

// 3. Создание мока через API (для скриптов)
if r.Method == http.MethodPost && r.URL.Path == "/_create" {
var req MockRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Bad Request", http.StatusBadRequest)
return
}
query := `INSERT INTO mock_endpoints (path, method, response_body, status_code) VALUES (?, ?, ?, ?)`
db.Exec(query, req.Path, req.Method, req.ResponseBody, req.StatusCode)
w.WriteHeader(http.StatusCreated)
fmt.Fprint(w, "✅ Мок создан")
return
}

// 4. Обработка самих моков
var responseBody string
var statusCode int
var contentType string
query := "SELECT response_body, status_code, content_type FROM mock_endpoints WHERE path = ? AND method = ?"
err := db.QueryRow(query, r.URL.Path, r.Method).Scan(&responseBody, &statusCode, &contentType)

if err == sql.ErrNoRows {
w.WriteHeader(http.StatusNotFound)
fmt.Fprintf(w, "Мок для пути %s не найден", r.URL.Path)
return
}

w.Header().Set("Content-Type", contentType)
w.WriteHeader(statusCode)
w.Write([]byte(responseBody))
})

log.Fatal(http.ListenAndServe(port, nil))
}

func renderDashboard(w http.ResponseWriter, db *sql.DB) {
rows, _ := db.Query("SELECT path, method, status_code FROM mock_endpoints ORDER BY id DESC")
defer rows.Close()

fmt.Fprint(w, `<html><head><title>MockWhale Dashboard</title><style>
body{font-family:'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;padding:40px;background:#f0f2f5;color:#1a1a1a}
.container{max-width:900px;margin:0 auto;background:#fff;padding:40px;border-radius:16px;box-shadow:0 10px 30px rgba(0,0,0,0.05)}
h1{color:#1f6feb;margin-bottom:30px;font-size:32px;display:flex;align-items:center;gap:15px}
table{width:100%;border-collapse:collapse;margin-top:20px}
th,td{padding:16px;border-bottom:1px solid #eee;text-align:left}
th{background:#fafafa;font-weight:600;color:#666;text-transform:uppercase;font-size:12px;letter-spacing:1px}
.btn{padding:10px 20px;border-radius:8px;text-decoration:none;font-size:14px;cursor:pointer;border:none;transition:0.2s;font-weight:500}
.btn-open{background:#e7f3ff;color:#1f6feb}
.btn-open:hover{background:#d0e7ff}
.form-card{background:#f8f9fa;padding:25px;border-radius:12px;margin-bottom:30px}
.form-group{display:flex;gap:12px}
input,select{padding:12px;border:2px solid #e1e4e8;border-radius:8px;font-size:14px;outline:none}
input:focus{border-color:#1f6feb}
.btn-add{background:#1f6feb;color:#fff}
.btn-add:hover{background:#1859be}
.method-tag{padding:4px 8px;border-radius:4px;font-size:12px;font-weight:bold;background:#eee}
</style></head><body><div class="container">`)

fmt.Fprint(w, "<h1>🐳 MockWhale</h1>")

fmt.Fprint(w, `<div class="form-card">
<h3 style="margin-top:0">Quick Create</h3>
<form action="/_create_web" method="POST" class="form-group">
<input name="path" placeholder="/api/v1/new-endpoint" required style="flex:2">
<select name="method" style="flex:0.5"><option>GET</option><option>POST</option><option>PUT</option><option>DELETE</option></select>
<input name="status" type="number" value="200" style="width:100px">
<button type="submit" class="btn btn-add">Add Mock</button>
</form>
</div>`)

fmt.Fprint(w, "<table><tr><th>Method</th><th>Endpoint Path</th><th>Status</th><th>Actions</th></tr>")
for rows.Next() {
var path, method string
var status int
rows.Scan(&path, &method, &status)
fmt.Fprintf(w, "<tr><td><span class='method-tag'>%s</span></td><td><code>%s</code></td><td>%d</td><td><a href='%s' class='btn btn-open' target='_blank'>Test API</a></td></tr>", method, path, status, path)
}
fmt.Fprint(w, "</table></div></body></html>")
}

func seedMock(db *sql.DB) {
var count int
db.QueryRow("SELECT COUNT(*) FROM mock_endpoints").Scan(&count)
if count == 0 {
db.Exec(`INSERT INTO mock_endpoints (path, method, response_body) VALUES (?, ?, ?)`, 
"/api/test", "GET", `{"message": "Hello from MockWhale, Chavs!", "status": "success"}`)
}
}

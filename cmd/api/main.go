package main

import (
"database/sql"
"encoding/json"
"fmt"
"log"
"net/http"
"github.com/Chavs23/mockwhale/internal/database"
)

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
if r.Method == http.MethodGet && r.URL.Path == "/_dashboard" {
renderDashboard(w, db)
return
}

if r.Method == http.MethodPost && r.URL.Path == "/_create_web" {
path := r.FormValue("path")
method := r.FormValue("method")
status := r.FormValue("status")
resp := r.FormValue("response")
db.Exec(`INSERT INTO mock_endpoints (path, method, response_body, status_code) VALUES (?, ?, ?, ?)`, 
path, method, resp, status)
http.Redirect(w, r, "/_dashboard", http.StatusSeeOther)
return
}

if r.Method == http.MethodPost && r.URL.Path == "/_delete" {
id := r.FormValue("id")
db.Exec(`DELETE FROM mock_endpoints WHERE id = ?`, id)
http.Redirect(w, r, "/_dashboard", http.StatusSeeOther)
return
}

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
rows, _ := db.Query("SELECT id, path, method, status_code FROM mock_endpoints ORDER BY id DESC")
defer rows.Close()

fmt.Fprint(w, `<html><head><title>MockWhale</title><style>
body{font-family:'Segoe UI',sans-serif;padding:40px;background:#f0f2f5;color:#1a1a1a}
.container{max-width:900px;margin:0 auto;background:#fff;padding:40px;border-radius:16px;box-shadow:0 10px 30px rgba(0,0,0,0.05)}
h1{color:#1f6feb;font-size:32px;margin-bottom:30px}
.form-card{background:#f8f9fa;padding:25px;border-radius:12px;margin-bottom:30px;border:1px solid #e1e4e8}
input, select, textarea{width:100%;padding:12px;margin:8px 0;border:2px solid #e1e4e8;border-radius:8px;font-family:monospace;outline:none}
input:focus, textarea:focus{border-color:#1f6feb}
.btn{padding:12px 24px;border-radius:8px;text-decoration:none;font-size:14px;cursor:pointer;border:none;font-weight:600;transition:0.2s}
.btn-add{background:#1f6feb;color:#fff;width:100%}
.btn-add:hover{background:#1859be}
.btn-del{background:#ff4d4d;color:#fff;padding:6px 12px;font-size:12px}
.btn-del:hover{background:#e60000}
table{width:100%;border-collapse:collapse;margin-top:20px}
th,td{padding:16px;border-bottom:1px solid #eee;text-align:left}
th{color:#666;font-size:12px;text-transform:uppercase;letter-spacing:1px}
</style></head><body><div class="container">`)

fmt.Fprint(w, `<h1>🐳 MockWhale</h1>
<div class="form-card">
<h3 style="margin-top:0">Create New Endpoint</h3>
<form action="/_create_web" method="POST">
<div style="display:flex;gap:12px">
<input name="path" placeholder="/api/v1/resource" required style="flex:2">
<select name="method" style="flex:0.5"><option>GET</option><option>POST</option><option>PUT</option><option>DELETE</option></select>
<input name="status" type="number" value="200" style="width:100px">
</div>
<textarea name="response" rows="4" placeholder='{"status": "ok"}'></textarea>
<button type="submit" class="btn btn-add">Add Mock Endpoint</button>
</form>
</div>
<table><tr><th>Method</th><th>Path</th><th>Status</th><th>Actions</th></tr>`)

for rows.Next() {
var id int
var path, method string
var status int
rows.Scan(&id, &path, &method, &status)
fmt.Fprintf(w, "<tr><td><b style='color:#1f6feb'>%s</b></td><td><code>%s</code></td><td>%d</td><td>
<div style='display:flex;gap:8px;align-items:center'>
<a href='%s' target='_blank' style='color:#1f6feb;font-size:14px'>Test</a>
<form action='/_delete' method='POST' style='margin:0'>
<input type='hidden' name='id' value='%d'>
<button class='btn btn-del'>Delete</button>
</form>
</div>
</td></tr>", method, path, status, path, id)
}
fmt.Fprint(w, "</table></div></body></html>")
}

func seedMock(db *sql.DB) {
var count int
db.QueryRow("SELECT COUNT(*) FROM mock_endpoints").Scan(&count)
if count == 0 {
db.Exec(`INSERT INTO mock_endpoints (path, method, response_body) VALUES (?, ?, ?)`, 
"/api/test", "GET", `{"message": "Hello from MockWhale, Chavs!"}`)
}
}

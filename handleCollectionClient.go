package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func GetCollectionClient(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM collection where status = 'active'")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var collections []Collection

	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.Slug); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		collections = append(collections, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Chuyển đổi dữ liệu sang định dạng JSON
	jsonData, err := json.Marshal(collections)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

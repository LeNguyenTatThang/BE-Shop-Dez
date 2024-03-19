package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func CreateCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var collection Collection
		err := json.NewDecoder(r.Body).Decode(&collection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		collection.Slug = createSlug(collection.Name)
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		result, err := db.Exec("INSERT INTO collection (name, status, slug) VALUES (?, ?, ?)", collection.Name, collection.Status, collection.Slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		collection.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collection)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func GetCollection(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM collection")
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
func GetDetailCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionID := r.URL.Query().Get("idcollection")
	if collectionID == "" {
		http.Error(w, "COllection ID is required", http.StatusBadRequest)
		return
	}

	// Mở kết nối đến cơ sở dữ liệu
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Thực hiện truy vấn để lấy thông tin chi tiết của thể loại
	var collectionName string
	var status string
	err = db.QueryRow("SELECT name, status FROM collection WHERE idcollection = ?", collectionID).Scan(&collectionName, &status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Phản hồi với thông tin chi tiết của thể loại
	response := map[string]string{"idcollection": collectionID, "name": collectionName, "status": status}
	json.NewEncoder(w).Encode(response)
}
func UpdateCollection(w http.ResponseWriter, r *http.Request) {
	collectionID := r.URL.Query().Get("id")
	if collectionID == "" {
		http.Error(w, "Missing collection ID", http.StatusBadRequest)
		return
	}
	var updatedCollection Collection
	err := json.NewDecoder(r.Body).Decode(&updatedCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	updatedCollection.Slug = createSlug(updatedCollection.Name)
	query := "UPDATE collection SET name = ?, status = ?, slug = ? WHERE idcollection = ?"
	_, err = db.Exec(query, updatedCollection.Name, updatedCollection.Status, updatedCollection.Slug, collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Collection updated successfully"))
}
func DeleteCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionID := r.URL.Query().Get("idcollection")
	if collectionID == "" {
		http.Error(w, "Collection ID is required", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM collection WHERE idcollection = ?", collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Collection deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}

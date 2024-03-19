package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func getProducts(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.IDProduct, &p.Description, &p.Thumbnail, &p.Content, &p.Status, &p.Slug, &p.Type, &p.IDCategory, &p.IDCollection, &p.CateIDCategory, &p.CollecIDCollection); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Chuyển đổi dữ liệu sang định dạng JSON
	jsonData, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Đặt header Content-Type là application/json
	w.Header().Set("Content-Type", "application/json")
	// Trả về dữ liệu JSON
	w.Write(jsonData)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// Xử lý logic tạo sản phẩm ở đây
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	// Xử lý logic cập nhật sản phẩm ở đây
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	// Xử lý logic xóa sản phẩm ở đây
}

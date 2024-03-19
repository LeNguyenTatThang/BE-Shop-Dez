package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// GetCategories trả về danh sách tất cả các category từ database
func GetCategories(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM category")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.Slug); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Chuyển đổi dữ liệu sang định dạng JSON
	jsonData, err := json.Marshal(categories)
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

// CreateCategory tạo một category mới trong database
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var category Category
		err := json.NewDecoder(r.Body).Decode(&category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Tạo giá trị cho trường Slug từ trường Name
		category.Slug = createSlug(category.Name)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		result, err := db.Exec("INSERT INTO category (name, status, slug) VALUES (?, ?, ?)", category.Name, category.Status, category.Slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		category.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Lấy chi tiết thể loại cần update
func GetCategoryDetails(w http.ResponseWriter, r *http.Request) {
	// Đảm bảo rằng yêu cầu chỉ được phép là GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lấy categoryID từ tham số URL
	categoryID := r.URL.Query().Get("idcategory")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
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
	var categoryName string
	var status string
	err = db.QueryRow("SELECT name, status FROM category WHERE idcategory = ?", categoryID).Scan(&categoryName, &status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Phản hồi với thông tin chi tiết của thể loại
	response := map[string]string{"idcategory": categoryID, "name": categoryName, "status": status}
	json.NewEncoder(w).Encode(response)
}

// UpdateCategory cập nhật một category trong database
func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	// Parse category ID from request URL or request body
	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		http.Error(w, "Missing category ID", http.StatusBadRequest)
		return
	}

	// Parse new category data from request body
	var updatedCategory Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Open database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	updatedCategory.Slug = createSlug(updatedCategory.Name)
	// Prepare SQL statement to update category
	query := "UPDATE category SET name = ?, status = ?, slug = ? WHERE idcategory = ?"
	_, err = db.Exec(query, updatedCategory.Name, updatedCategory.Status, updatedCategory.Slug, categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Category updated successfully"))
}

// DeleteCategory xóa một category từ database
func DeleteCategory(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Đảm bảo rằng yêu cầu có tham số idcategory
	categoryID := r.URL.Query().Get("idcategory")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	// Mở kết nối đến cơ sở dữ liệu
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopthoitrang")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Thực hiện truy vấn xóa
	_, err = db.Exec("DELETE FROM category WHERE idcategory = ?", categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Thêm tiêu đề CORS vào phản hồi khi xóa thành công
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Category deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}

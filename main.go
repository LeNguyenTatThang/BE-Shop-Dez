package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gorilla/handlers"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	_ "github.com/go-sql-driver/mysql"
)

// Product là struct để ánh xạ dữ liệu từ bảng "product"
type Product struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	IDProduct          string `json:"id_product"`
	Description        string `json:"description"`
	Thumbnail          string `json:"thumbnail"`
	Content            string `json:"content"`
	Status             string `json:"status"`
	Slug               string `json:"slug"`
	Type               string `json:"type"`
	IDCategory         int    `json:"id_category"`
	IDCollection       int    `json:"id_collection"`
	CateIDCategory     int    `json:"cate_id_category"`
	CollecIDCollection int    `json:"collec_id_collection"`
}

type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Slug   string `json:"slug"`
}

func main() {
	//product
	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/products/create", createProduct)
	http.HandleFunc("/products/update", updateProduct)
	http.HandleFunc("/products/delete", deleteProduct)
	//category
	http.HandleFunc("/categorys", GetCategories)
	http.HandleFunc("/create-category", CreateCategory)
	http.HandleFunc("/detail-category", GetCategoryDetails)
	http.HandleFunc("/update-category", UpdateCategory)
	http.HandleFunc("/delete-category", DeleteCategory)

	// Tạo middleware CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Áp dụng middleware CORS cho router chính
	log.Fatal(http.ListenAndServe(":8080", cors(http.DefaultServeMux)))
}

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

func createSlug(name string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	processedString, _, e := transform.String(t, name)
	if e != nil {
		panic(e)
	}

	processedString = strings.ReplaceAll(processedString, " ", "-")
	processedString = strings.ToLower(processedString)
	processedString = strings.Trim(processedString, "-")
	processedString = strings.ReplaceAll(processedString, "--", "-")

	return processedString
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
	err = db.QueryRow("SELECT name FROM category WHERE idcategory = ?", categoryID).Scan(&categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Phản hồi với thông tin chi tiết của thể loại
	response := map[string]string{"idcategory": categoryID, "name": categoryName}
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

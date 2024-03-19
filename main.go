package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/gorilla/handlers"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	_ "github.com/go-sql-driver/mysql"
)

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
	//collection
	http.HandleFunc("/collection", GetCollection)
	http.HandleFunc("/create-collection", CreateCollection)
	http.HandleFunc("/detail-collection", GetDetailCollection)
	http.HandleFunc("/update-collection", UpdateCollection)
	http.HandleFunc("/delete-collection", DeleteCollection)
	//upload
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/get-image-url", getImageURL)
	fs := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs))
	//client
	http.HandleFunc("/collection-client", GetCollectionClient)
	// Tạo middleware CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	// Áp dụng middleware CORS cho router chính
	log.Fatal(http.ListenAndServe(":8080", cors(http.DefaultServeMux)))
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		log.Println("Error retrieving the file", err)
		return
	}
	defer file.Close()
	if handler.Header.Get("Content-Type") != "image/jpeg" && handler.Header.Get("Content-Type") != "image/png" {
		http.Error(w, "Unsupportd file format. Please choose img with jpeg/png", http.StatusBadRequest)
		return
	}
	f, err := os.OpenFile("img/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	if _, err := io.Copy(f, file); err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		log.Println("Error copying the file:", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

func getImageURL(w http.ResponseWriter, r *http.Request) {
	// Đọc tên file hình ảnh từ tham số URL
	fileName := r.URL.Query().Get("filename")

	// Kiểm tra xem tên file có tồn tại không
	if fileName == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	// Tạo URL cho hình ảnh
	imageURL := "http://localhost:8080/img/" + fileName

	// Tạo phản hồi JSON chứa đường dẫn hình ảnh
	response := ImageResponse{ImageURL: imageURL}

	// Thiết lập header Content-Type là application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode phản hồi JSON và gửi nó về client
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

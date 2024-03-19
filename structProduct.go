package main

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

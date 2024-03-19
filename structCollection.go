package main

type Collection struct {
	ID     int    `json: "idcollection"`
	Name   string `json: "name"`
	Status string `status: "status"`
	Slug   string `json: "slug"`
}

package main

type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Slug   string `json:"slug"`
}

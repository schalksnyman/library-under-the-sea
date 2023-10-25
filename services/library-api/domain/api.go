package domain

import "net/http"

type Client interface {
	Add(res http.ResponseWriter, req *http.Request)
	FindAll(res http.ResponseWriter, req *http.Request)
	FindByTitle(res http.ResponseWriter, req *http.Request)
	FindBook(res http.ResponseWriter, req *http.Request)
	UpdateBookTitle(res http.ResponseWriter, req *http.Request)
	DeleteBook(res http.ResponseWriter, req *http.Request)
}

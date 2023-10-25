package logical

import (
	"encoding/json"
	libraryAPI "library-under-the-sea/services/library-api/domain"
	library "library-under-the-sea/services/library/domain"
	"net/http"
)

var _ libraryAPI.Client = (*client)(nil)

func New(l library.Client) *client {

	return &client{
		l: l,
	}
}

type client struct {
	l library.Client
}

func (c *client) Add(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var book library.Book
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: "Invalid Payload"})
		return
	}
	id, err := c.l.SaveBook(book)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
	result := make(map[string]string)
	result["id"] = id
	json.NewEncoder(res).Encode(result)
}

func (c *client) FindAll(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	results, err := c.l.ListAll()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(results)
}

func (c *client) FindByTitle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var request libraryAPI.FindByTitleRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: "Invalid Payload"})
		return
	}

	results, err := c.l.ListBooksByTitle(request.Title)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(results)
}

func (c *client) FindBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var request libraryAPI.FindBookRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: "Invalid Payload"})
		return
	}

	result, err := c.l.FindBook(request.Id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(result)
}

func (c *client) UpdateBookTitle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var request libraryAPI.UpdateTitleRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: "Invalid Payload"})
		return
	}

	err = c.l.UpdateTitle(request.Id, request.Title)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (c *client) DeleteBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var request libraryAPI.DeleteBookRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: "Invalid Payload"})
		return
	}

	err = c.l.DeleteBook(request.Id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(libraryAPI.ErrorResponse{Message: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
}

package transport

import (
	"encoding/json"
	"net/http"

	"youtubeviews/models"
	"youtubeviews/service"
)

type HttpTransport struct {
	service service.Service
}

func NewHttpTransport(service service.Service) *HttpTransport {
	return &HttpTransport{service: service}
}

func (t *HttpTransport) Increment(writer http.ResponseWriter, request *http.Request) {
	// Declare a variable to hold the request payload
	var req models.IncrementPayload

	// Check if the request body is empty
	if request.Body == nil {
		http.Error(writer, "empty request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	// Decode the request body into the payload struct
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to increment the view count
	views, inc, err := t.service.Increment(request.Context(), req.VideoID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON and encode the response
	incrementResponse := models.IncrementViewResponse{Views: views, Increment: inc}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(incrementResponse)
	writer.WriteHeader(http.StatusOK)
}

func (t *HttpTransport) Get(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "empty request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()
	var req models.ViewCountPayload
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	views, err := t.service.Get(request.Context(), req.VideoID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	getResponse := models.ViewCountResponse{Views: views}
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(getResponse); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}
}

func (t *HttpTransport) GetTopVideos(writer http.ResponseWriter, request *http.Request) {
	// Check if the request body is empty
	if request.Body == nil {
		http.Error(writer, "empty request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	// Declare a variable to hold the request payload
	var req models.GetTopVideosPayload

	// Decode the request body into the payload struct
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to get the top videos
	response, err := t.service.GetTopVideos(request.Context(), req.Page, req.Limit)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON and encode the response
	topVideosResponse := models.GetTopVideosResponse{TopVideos: response}
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(topVideosResponse); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}
}

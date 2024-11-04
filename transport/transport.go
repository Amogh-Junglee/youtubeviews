package transport

import (
	"encoding/json"
	"net/http"

	models "YoutubeViews/models"
	service "YoutubeViews/service"
)

type HttpTransport struct {
	service service.Service
}

func NewHttpTransport(service service.Service) *HttpTransport {
	return &HttpTransport{service: service}
}

func (t *HttpTransport) Increment(writer http.ResponseWriter, request *http.Request) {
	var req models.IncrementPayload
	if request.Body == nil {
		http.Error(writer, "empty request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.VideoID == "" {
		http.Error(writer, "videoId is required", http.StatusBadRequest)
		return
	}

	response, err := t.service.Increment(request.Context(), req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
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

	if req.VideoID == "" {
		http.Error(writer, "videoId is required", http.StatusBadRequest)
		return
	}

	response, err := t.service.Get(request.Context(), req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}
}

func (t *HttpTransport) GetTopVideos(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "empty request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()
	var req models.GetTopVideosPayload
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := t.service.GetTopVideos(request.Context(), req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}
}

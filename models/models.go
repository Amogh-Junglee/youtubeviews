package models

type IncrementPayload struct {
	VideoID string `json:"videoId"`
}

type IncrementViewResponse struct {
	Views     int `json:"views"`
	Increment int `json:"increment"`
}

type ViewCountPayload struct {
	VideoID string `json:"videoId"`
}

type ViewCountResponse struct {
	Views int `json:"views"`
}

func (r ViewCountResponse) GetViews() int {
	return r.Views
}

type GetTopVideosPayload struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetTopVideosResponse struct {
	TopVideos []map[string]interface{} `json:"top_videos"`
}

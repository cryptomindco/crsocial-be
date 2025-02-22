package portal

type PostContentRequest struct {
	Content string `json:"postContent"`
}

type PostRequest struct {
	Id uint64 `json:"id"`
}

type LikeRequest struct {
	Id     uint64 `json:"id"`
	Status bool   `json:"status"`
}

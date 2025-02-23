package portal

type PostContentRequest struct {
	Content string `json:"postContent"`
}

type PostRequest struct {
	Id uint64 `json:"id"`
}

type LikeTargetRequest struct {
	Id   uint64 `json:"id"`
	Type int    `json:"type"`
}

type ReplyRequest struct {
	TargetId uint64 `json:"targetId"`
	Reply    string `json:"reply"`
	Type     int    `json:"type"`
}

type LikeRequest struct {
	Id     uint64 `json:"id"`
	Status bool   `json:"status"`
}

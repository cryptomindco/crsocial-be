package portal

type GetTimelineRequest struct {
	Cursor string `validate:"omitempty,cursor"`
	Limit  int64  `validate:"omitempty,limit"`
}
type ImagesURL struct {
	Urls []string `validate:"omitempty,urls"`
	Url  string   `validate:"omitempty,url"`
}

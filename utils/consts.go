package utils

const (
	SortByCreated  = 1
	SortByLastSeen = 2

	SortASC                     = 1
	SortDESC                    = 2
	LimitOfFetchTimeline        = 50
	ImageType            string = "image"
	VideoType            string = "video"
	OtherType            string = "other"
	FileSeparate         string = ";;"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")

type ResponseData struct {
	IsError bool        `json:"error"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type UserRole int
type FollowStatus uint
type LikeType int

const (
	UserRoleNone UserRole = iota
	UserRoleAdmin
)

const (
	PostType LikeType = iota
	CommentType
)

const (
	UnFollow FollowStatus = iota
	Follow
)

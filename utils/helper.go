package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func NewError(err error, code int) *Error {
	return &Error{
		error: err,
		Code:  code,
	}
}
func ResponseOK(w http.ResponseWriter, data interface{}, errs ...*Error) {
	if len(errs) > 0 && errs[0] != nil {
		Response(w, http.StatusOK, errs[0], data)
		return
	}
	Response(w, http.StatusOK, nil, data)
}

func Response(w http.ResponseWriter, httpStatus int, err error, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	res := response{
		Success: err == nil,
		Code:    StatusOK,
		Message: "ok",
		Data:    data,
	}

	if err != nil {
		switch er := err.(type) {
		case *Error:
			res.Code = er.Code
			res.Message = er.Error()
			w.WriteHeader(er.HttpStatus())
		default:
			res.Message = err.Error()
			res.Code = ErrorInternalCode
			w.WriteHeader(httpStatus)
		}
	} else {
		w.WriteHeader(httpStatus)
	}
	enc.Encode(res)
}

func IsEmpty(x interface{}) bool {
	switch value := x.(type) {
	case string:
		return value == ""
	case int32:
		return value == 0
	case int:
		return value == 0
	case uint32:
		return value == 0
	case uint64:
		return value == 0
	case int64:
		return value == 0
	case float64:
		return value == 0
	case bool:
		return false
	default:
		return true
	}
}

func DecodeQuery(object interface{}, query url.Values) error {
	err := decoder.Decode(object, query)
	if err != nil {
		return err
	}

	return nil
}

func SetValue[T any](source *T, value T) {
	if !IsEmpty(value) && source != &value {
		*source = value
	}
}

func ImageToBase64(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return "", err
	}

	qrBytesString := buf.Bytes()
	imgBase64Str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrBytesString)

	return imgBase64Str, nil
}

var imagePath = getBinPath() + "/upload/images"
var videoPath = getBinPath() + "/upload/videos"
var otherPath = getBinPath() + "/upload/other"

func getBinPath() string {
	e, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return e
}

func GetImagePath() string {
	return imagePath
}

func GetVideoPath() string {
	return videoPath
}

func GetFileDirectory(fileType string) string {
	switch fileType {
	case "video":
		return GetVideoPath()
	case "image":
		return GetImagePath()
	default:
		return otherPath
	}
}

func GetConvertFilename(uploadDir string, fileName string) string {
	fileExt := GetFileExtension(fileName)
	if IsEmpty(fileExt) {
		return ""
	}
	newName := ""
	for {
		randName := RandSeq(48)
		newName = fmt.Sprintf("%s.%s", randName, fileExt)
		fileExist := CheckFileExist(uploadDir, newName)
		if !fileExist {
			break
		}
	}
	return newName
}

func CheckFileExist(dir, fileName string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", dir, fileName))
	return !os.IsNotExist(err)
}

func GetFileExtension(fileName string) string {
	nameArr := strings.Split(fileName, ".")
	if len(nameArr) < 2 {
		return ""
	}
	return nameArr[len(nameArr)-1]
}

func GetFileType(originalType string) string {
	if strings.HasPrefix(originalType, "video/") {
		return VideoType
	}
	if strings.HasPrefix(originalType, "image/") {
		return ImageType
	}
	return OtherType
}

func ConvertUrlArray(siteRoot string, urlArray []string) []string {
	res := make([]string, 0)
	for _, url := range urlArray {
		siteUrl := ConvertToSiteUrl(siteRoot, url)
		res = append(res, siteUrl)
	}
	return res
}

func ConvertToSiteUrl(siteRoot, fileUrl string) string {
	uploadPrefix := "/upload"
	pos := strings.Index(fileUrl, uploadPrefix)
	if pos > -1 {
		fileUrl = fileUrl[pos:]
	} else {
		return "#"
	}
	result := fmt.Sprintf("%s%s", siteRoot, fileUrl)
	return result
}

func IsVideoFileType(fileMainType string) bool {
	return fileMainType == VideoType
}

func IsImageFileType(fileMainType string) bool {
	return fileMainType == ImageType
}

func ConvertImageToBase64(fileName string) string {
	imgFile, err := os.Open(imagePath + "/" + fileName) //Image file

	if err != nil {
		log.Println(err)
		return ""
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	return imgBase64Str
}

func ContainsUint64(inputArr []uint64, value uint64) bool {
	if inputArr == nil {
		return false
	}
	for _, v := range inputArr {
		if v == value {
			return true
		}
	}
	return false
}

// convert time to string, format yyyy-mm-dd hh:MM:ss
func TimeToStringWithoutTimeZone(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func GetUserDisplayName(userName string, displayName string) string {
	if IsEmpty(displayName) {
		return userName
	}
	return displayName
}

// handler date format (YYYY/MM/DD)
func HandlerDateFormat(date string) string {
	dateArr := strings.Split(date, "/")
	if len(dateArr) < 3 {
		return date
	}

	if utf8.RuneCountInString(dateArr[1]) == 1 {
		dateArr[1] = fmt.Sprintf("0%s", dateArr[1])
	}

	if utf8.RuneCountInString(dateArr[2]) == 1 {
		dateArr[2] = fmt.Sprintf("0%s", dateArr[2])
	}
	return strings.Join(dateArr, "/")
}

func GetSecondDurationFromStartEnd(startTime, endTime time.Time) uint64 {
	return uint64(math.Floor(endTime.Sub(startTime).Seconds()))
}

func CatchObject(from interface{}, to interface{}) error {
	jsonBytes, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, &to)
	if err != nil {
		return err
	}
	return nil
}

func RequestBodyToString(body io.ReadCloser) string {
	b, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(b)
}

func ObjectToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func JsonStringToObject(jsonString string, to interface{}) error {
	err := json.Unmarshal([]byte(jsonString), &to)
	if err != nil {
		return err
	}
	return nil
}

// handler error from rpc
func HandlerRPCError(err error) error {
	if err == nil {
		return nil
	}
	errString := err.Error()
	if strings.HasPrefix(errString, "rpc error:") {
		errSplitArr := strings.Split(errString, "desc =")
		if len(errSplitArr) < 2 {
			return err
		}
		errMsg := strings.TrimSpace(errSplitArr[1])
		return fmt.Errorf("%s", errMsg)
	}
	return err
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetHandleFromUsername(server, username string) string {
	serverName := server
	if strings.HasPrefix(server, "http") {
		serverArr := strings.Split(server, "//")
		if len(serverArr) > 1 {
			serverName = serverArr[1]
		}
	}
	return fmt.Sprintf("%s.%s", username, serverName)
}

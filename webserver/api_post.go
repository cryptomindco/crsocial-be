package webserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"socialat/be/storage"
	"socialat/be/utils"
	"socialat/be/webserver/portal"
	"strconv"
	"time"
)

type apiPost struct {
	*WebServer
}

func (a *apiPost) PostWithFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	var fileNumber = len(r.MultipartForm.File)
	var content = r.Form.Get("postContent")
	var filesType = r.Form.Get("fileType")
	uploadDir := utils.GetFileDirectory(utils.GetFileType(filesType))
	fileArr := make([]string, 0)
	claims, _ := a.credentialsInfo(r)
	for i := 0; i < fileNumber; i++ {
		file, handler, err := r.FormFile("files[" + strconv.Itoa(i) + "]")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		defer file.Close()
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("Create folder failed")
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		newFileUpload := uploadDir + "/" + utils.GetConvertFilename(uploadDir, handler.Filename)
		err = ioutil.WriteFile(newFileUpload, fileBytes, 0777)
		if err != nil {
			fmt.Println("Write file error")
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		fileArr = append(fileArr, utils.ConvertToSiteUrl(a.conf.SiteRoot, newFileUpload))
	}
	var postErr error
	if len(fileArr) > 0 {
		_, postErr = a.service.CreateNewPost(claims.Id, claims.UserName, content, filesType, fileArr)
	} else {
		_, postErr = a.service.CreateNewPost(claims.Id, claims.UserName, content, "", nil)
	}
	if postErr != nil {
		utils.Response(w, http.StatusInternalServerError, postErr, nil)
		return
	}
	utils.ResponseOK(w, Map{})
}

func (a *apiPost) likePost(w http.ResponseWriter, r *http.Request) {
	var f portal.PostRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	isLiked, err := a.service.CheckLoginLiked(claims.UserName, f.Id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	var post storage.Post
	err = a.db.GetById(f.Id, &post)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	if !isLiked {
		err = a.db.CreateLike(&storage.Like{
			Username: claims.UserName,
			PostId:   f.Id,
			LikedAt:  time.Now(),
		})
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		post.LikeCount++
	} else {
		tx := a.db.GetDB().Begin()
		if err := tx.Where("username = ? AND post_id = ?", claims.UserName, f.Id).Delete(storage.Like{}).Error; err != nil {
			tx.Rollback()
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		tx.Commit()
		post.LikeCount--
	}
	// update post
	err = a.db.UpdatePost(&post)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, !isLiked)
}

func (a *apiPost) PostWithoutFiles(w http.ResponseWriter, r *http.Request) {
	var f portal.PostContentRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	_, postErr := a.service.CreateNewPost(claims.Id, claims.UserName, f.Content, "", nil)
	if postErr != nil {
		utils.Response(w, http.StatusInternalServerError, postErr, nil)
		return
	}
	utils.ResponseOK(w, Map{})
}

func (a *apiPost) uploadImages(w http.ResponseWriter, r *http.Request) {
	// var imgUrls portal.ImagesURL
	// err := a.parseJSONAndValidate(r, &imgUrls)
	// if err != nil {
	// 	log.Errorf("get image url params failed: %v", err)
	// 	utils.Response(w, http.StatusInternalServerError, err, nil)
	// 	return
	// }
	// if len(imgUrls.Urls) == 0 && utils.IsEmpty(imgUrls.Url) {
	// 	log.Errorf("image from url is empty")
	// 	utils.Response(w, http.StatusInternalServerError, fmt.Errorf("image from url is empty"), nil)
	// 	return
	// }
	// images := []atlib.Image{}
	// if len(imgUrls.Urls) > 0 {
	// 	for _, urlStr := range imgUrls.Urls {
	// 		urlObj, err := url.Parse(urlStr)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		images = append(images, atlib.Image{
	// 			Title: "Image title",
	// 			Uri:   *urlObj,
	// 		})
	// 	}
	// } else {
	// 	urlObj, err := url.Parse(imgUrls.Url)
	// 	if err == nil {
	// 		images = append(images, atlib.Image{
	// 			Title: "Image title",
	// 			Uri:   *urlObj,
	// 		})
	// 	}
	// }
	// if len(images) == 0 {
	// 	log.Errorf("parse image data failed")
	// 	utils.Response(w, http.StatusInternalServerError, fmt.Errorf("parse image data failed"), nil)
	// 	return
	// }
	// claims, _ := a.credentialsInfo(r)
	// utils.ResponseOK(w, blobs)
}

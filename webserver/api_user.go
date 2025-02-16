package webserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"socialat/be/storage"
	"socialat/be/utils"
	"socialat/be/webserver/portal"
	"time"

	"github.com/go-chi/chi/v5"
)

type apiUser struct {
	*WebServer
}

func (a *apiUser) getUserPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	var author storage.Author
	if username == "-1" {
		claims, _ := a.credentialsInfo(r)
		author = storage.Author{
			Username: claims.UserName,
			FullName: claims.FullName,
			Avatar:   claims.Avatar,
		}
		username = claims.UserName
	} else {
		var err error
		// get user by id
		userInfo, err := a.service.GetUserInfoByUsername(username)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		author = storage.Author{
			Username: userInfo.Username,
			FullName: userInfo.FullName,
			Avatar:   userInfo.Avatar,
		}
	}
	posts, err := a.service.GetPostListByUsername(username)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, utils.NewError(err, utils.ErrorInternalCode), nil)
		return
	}
	postViews := make([]storage.PostView, 0)
	for _, post := range posts {
		postView := storage.PostView{
			Post:   &post,
			Author: &author,
		}
		postViews = append(postViews, postView)
	}
	postViews = storage.HandlerPostFileUrls(a.conf.SiteRoot, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiUser) getAllPosts(w http.ResponseWriter, r *http.Request) {
	var posts = make([]storage.Post, 0)
	posts, _ = a.service.GetAllPost(40)
	postViews := make([]storage.PostView, 0)
	authorMap := make(map[string]storage.Author)
	for _, post := range posts {
		author, exist := authorMap[post.Username]
		if !exist {
			userInfo, err := a.service.GetUserInfoByUsername(post.Username)
			if err == nil {
				author = storage.Author{
					Username: userInfo.Username,
					FullName: userInfo.FullName,
					Avatar:   userInfo.Avatar,
				}
				authorMap[post.Username] = author
			}
		}
		postView := storage.PostView{
			Post:   &post,
			Author: &author,
		}
		postViews = append(postViews, postView)
	}
	postViews = storage.HandlerPostFileUrls(a.conf.SiteRoot, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiUser) getTimelines(w http.ResponseWriter, r *http.Request) {
	claims, _ := a.credentialsInfo(r)
	var posts = make([]storage.Post, 0)
	if claims.IsLogin {
		posts, _ = a.service.GetPostOfFollowing(claims.UserName)
	}
	postViews := make([]storage.PostView, 0)
	authorMap := make(map[string]storage.Author)
	for _, post := range posts {
		author, exist := authorMap[post.Username]
		if !exist {
			userInfo, err := a.service.GetUserInfoByUsername(post.Username)
			if err == nil {
				author = storage.Author{
					Username: userInfo.Username,
					FullName: userInfo.FullName,
					Avatar:   userInfo.Avatar,
				}
				authorMap[post.Username] = author
			}
		}
		postView := storage.PostView{
			Post:   &post,
			Author: &author,
		}
		postViews = append(postViews, postView)
	}
	postViews = storage.HandlerPostFileUrls(a.conf.SiteRoot, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiUser) getUserByName(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userInfo, err := a.service.GetUserInfoByUsername(username)
	claims, _ := a.credentialsInfo(r)
	followStatus := false
	followings := make([]storage.FollowingDisplay, 0)
	if claims.IsLogin {
		if claims.UserName != username {
			follow, err := a.service.GetFollowByName(claims.UserName, username)
			if err == nil && follow != nil {
				followStatus = follow.Status
			}
		} else {
			following, err := a.service.GetFollowingOfUser(username)
			if err == nil {
				for _, targetFollow := range following {
					targetInfo, err := a.service.GetUserInfoByUsername(targetFollow.Target)
					if err == nil {
						followDisplay := storage.FollowingDisplay{
							Target:   targetInfo.Username,
							Avatar:   targetInfo.Avatar,
							FullName: targetInfo.FullName,
						}
						followings = append(followings, followDisplay)
					}
				}
			}
		}
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, Map{"userinfo": userInfo, "follow": followStatus, "following": followings})
}

func (a *apiUser) updateDisplayName(w http.ResponseWriter, r *http.Request) {
	var f portal.UpdateUserInfoRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	userInfo, err := a.service.GetUserInfoByUsername(claims.UserName)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	userInfo.FullName = f.FullName
	err = a.db.UpdateUserInfo(userInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, userInfo)
}

func (a *apiUser) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	defer file.Close()
	err = os.MkdirAll(utils.GetImagePath(), os.ModePerm)
	if err != nil {
		fmt.Println("Create folder failed")
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	uploadDir := utils.GetImagePath()
	newFileUpload := uploadDir + "/" + utils.GetConvertFilename(uploadDir, handler.Filename)
	err = ioutil.WriteFile(newFileUpload, fileBytes, 0777)
	if err != nil {
		fmt.Println("Write file error")
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	userInfo, err := a.service.GetUserInfoByUsername(claims.UserName)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	userInfo.Avatar = utils.ConvertToSiteUrl(a.conf.SiteRoot, newFileUpload)
	err = a.db.UpdateUserInfo(userInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, userInfo)
}

func (a *apiUser) updateBio(w http.ResponseWriter, r *http.Request) {
	var f portal.UpdateUserInfoRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	userInfo, err := a.service.GetUserInfoByUsername(claims.UserName)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	userInfo.Bio = f.Bio
	err = a.db.UpdateUserInfo(userInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, userInfo)
}

func (a *apiUser) FollowUpdateUser(w http.ResponseWriter, r *http.Request) {
	var f portal.FollowRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)

	follower, err := a.service.GetFollowByName(claims.UserName, f.TargetUsername)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	// follow: true, unfollow: false
	if follower == nil {
		follower = &storage.Follower{
			Follower:   claims.UserName,
			Target:     f.TargetUsername,
			FollowedAt: time.Now(),
			Status:     f.Status,
		}
		err = a.db.CreateFollower(follower)
	} else {
		follower.Status = f.Status
		err = a.db.UpdateFollower(follower)
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, follower.Status)
}

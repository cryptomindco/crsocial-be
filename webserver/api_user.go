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
)

type apiUser struct {
	*WebServer
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

func (a *apiUser) updateFullProfile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	var fullName = r.Form.Get("fullName")
	var bio = r.Form.Get("bio")
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
	userInfo.FullName = fullName
	userInfo.Bio = bio
	err = a.db.UpdateUserInfo(userInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, userInfo)
}

func (a *apiUser) updateProfileInfo(w http.ResponseWriter, r *http.Request) {
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
	userInfo.FullName = f.FullName
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

func (a *apiUser) reply(w http.ResponseWriter, r *http.Request) {
	var f portal.ReplyRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	if f.Type != int(utils.PostType) && f.Type != int(utils.CommentType) {
		utils.Response(w, http.StatusBadRequest, fmt.Errorf("like type error"), nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	comment := storage.Comment{
		Username:  claims.UserName,
		Content:   f.Reply,
		CreatedAt: time.Now(),
	}
	if f.Type == int(utils.PostType) {
		comment.PostId = f.TargetId
		// if reply to comment
		// get post
		var post storage.Post
		err = a.db.GetById(comment.PostId, &post)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		post.CommentCount++
		err := a.db.UpdatePost(&post)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
	} else {
		// get parent comment
		var parent storage.Comment
		err = a.db.GetById(f.TargetId, &parent)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		comment.PostId = parent.PostId
		comment.ParentId = parent.Id
	}
	// create comment
	err = a.db.CreateComment(&comment)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	commentView := storage.CommentView{
		Id: comment.Id,
		Author: &storage.Author{
			Username: claims.UserName,
			FullName: claims.FullName,
			Avatar:   claims.Avatar,
		},
		Content:   comment.Content,
		ParentId:  comment.ParentId,
		LikeCount: 0,
		CreatedAt: comment.CreatedAt,
	}
	utils.ResponseOK(w, commentView)
}

func (a *apiUser) likeHandle(w http.ResponseWriter, r *http.Request) {
	var f portal.LikeTargetRequest
	err := a.parseJSONAndValidate(r, &f)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	if f.Type != int(utils.PostType) && f.Type != int(utils.CommentType) {
		utils.Response(w, http.StatusBadRequest, fmt.Errorf("like type error"), nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	isLiked, err := a.service.CheckLoginLiked(claims.UserName, f.Id, f.Type)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	if !isLiked {
		err = a.db.CreateLike(&storage.Like{
			Username: claims.UserName,
			TargetId: f.Id,
			Type:     f.Type,
			LikedAt:  time.Now(),
		})
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
	} else {
		tx := a.db.GetDB().Begin()
		if err := tx.Where("username = ? AND target_id = ? AND type = ?", claims.UserName, f.Id, f.Type).Delete(storage.Like{}).Error; err != nil {
			tx.Rollback()
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
		}
		tx.Commit()
	}
	if f.Type == int(utils.PostType) {
		err = a.service.UpdatePostLikeCount(uint32(f.Id), !isLiked)
	} else {
		err = a.service.UpdateCommentLikeCount(uint32(f.Id), !isLiked)
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	utils.ResponseOK(w, !isLiked)
}

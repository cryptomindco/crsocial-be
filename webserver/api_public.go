package webserver

import (
	"net/http"
	"socialat/be/storage"
	"socialat/be/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type apiPublic struct {
	*WebServer
}

func (a *apiPublic) getUserPosts(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	claims, _ := a.credentialsInfo(r)
	var author *storage.Author
	if username == "-1" {
		author = &storage.Author{
			Username: claims.UserName,
			FullName: claims.FullName,
			Avatar:   claims.Avatar,
		}
		username = claims.UserName
	} else {
		var err error
		// get user by id
		author, err = a.service.GetAuthorByUsername(username)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, err, nil)
			return
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
			Author: author,
		}
		postViews = append(postViews, postView)
	}
	postViews = a.service.HandleForPostList(a.conf.SiteRoot, claims.UserName, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiPublic) getPostDetail(w http.ResponseWriter, r *http.Request) {
	postIdStr := chi.URLParam(r, "id")
	postId, err := strconv.ParseUint(postIdStr, 0, 32)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, err, nil)
		return
	}
	claims, _ := a.credentialsInfo(r)
	// Get Post
	var post storage.Post
	err = a.db.GetById(postId, &post)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	author, err := a.service.GetAuthorByUsername(post.Username)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	commentViews, err := a.service.GetCommentViewList(claims.UserName, post.Id)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, err, nil)
		return
	}
	postView := storage.PostView{
		Post:     &post,
		Author:   author,
		Comments: commentViews,
	}
	postView = a.service.HandleForPost(a.conf.SiteRoot, claims.UserName, postView)
	utils.ResponseOK(w, postView)
}

func (a *apiPublic) getAllPosts(w http.ResponseWriter, r *http.Request) {
	var posts = make([]storage.Post, 0)
	posts, _ = a.service.GetAllPost(40)
	claims, _ := a.credentialsInfo(r)
	postViews := make([]storage.PostView, 0)
	authorMap := make(map[string]*storage.Author)
	for _, post := range posts {
		author, exist := authorMap[post.Username]
		if !exist {
			var err error
			author, err = a.service.GetAuthorByUsername(post.Username)
			if err == nil {
				authorMap[post.Username] = author
			}
		}
		postView := storage.PostView{
			Post:   &post,
			Author: author,
		}
		postViews = append(postViews, postView)
	}
	postViews = a.service.HandleForPostList(a.conf.SiteRoot, claims.UserName, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiPublic) getTimelines(w http.ResponseWriter, r *http.Request) {
	claims, _ := a.credentialsInfo(r)
	var posts = make([]storage.Post, 0)
	if claims.IsLogin {
		posts, _ = a.service.GetPostOfFollowing(claims.UserName)
	}
	postViews := make([]storage.PostView, 0)
	authorMap := make(map[string]*storage.Author)
	for _, post := range posts {
		author, exist := authorMap[post.Username]
		if !exist {
			var err error
			author, err = a.service.GetAuthorByUsername(post.Username)
			if err == nil {
				authorMap[post.Username] = author
			}
		}
		// get comments
		postView := storage.PostView{
			Post:   &post,
			Author: author,
		}
		postViews = append(postViews, postView)
	}
	postViews = a.service.HandleForPostList(a.conf.SiteRoot, claims.UserName, postViews)
	utils.ResponseOK(w, Map{
		"posts": postViews,
	})
}

func (a *apiPublic) getUserByName(w http.ResponseWriter, r *http.Request) {
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

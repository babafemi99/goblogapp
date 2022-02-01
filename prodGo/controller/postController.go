package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prodGo/cache"
	"prodGo/entity"
	"prodGo/er"
	"prodGo/service"
	"strings"
	// "strings"
)

var (
	serviceRepo service.PostService
	postCache   cache.PostCache
)

type controller struct{}

type PostController interface {
	GetPost(rw http.ResponseWriter, r *http.Request)
	AddPost(rw http.ResponseWriter, r *http.Request)
	GetPostById(rw http.ResponseWriter, r *http.Request)
}

func NewPostController(service service.PostService, cache cache.PostCache) PostController {
	serviceRepo = service
	postCache = cache
	return &controller{}
}

func (*controller) GetPost(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	posts, err := serviceRepo.FindAll()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode((er.CError{Message: "Error getting Data- from here "}))
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(posts)
}

func (*controller) GetPostById(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("the post id in this case is")
	var post *entity.Post
	var err error
	rw.Header().Set("Content-type", "application/json")
	postId := strings.Split(r.URL.Path, "/")[2]
	fmt.Println("the post id in this case is", postId)
	post = postCache.Get(postId)
	if post == nil {
	post, err = serviceRepo.FIndById(postId)
	fmt.Println("no cache", post)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode((er.CError{Message: "Error getting Data- from here "}))
		return
	}
	postCache.Set(postId, post)
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(post)
}

func (*controller) AddPost(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode((er.CError{Message: "Error Decodiing Data"}))
		return
	}
	validateErr := serviceRepo.Validate(&post)
	if validateErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode((er.CError{Message: validateErr.Error()}))
		return
	}
	result, createErr := serviceRepo.Create(&post)
	if createErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode((er.CError{Message: "Error Adding Data"}))
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(result)
}

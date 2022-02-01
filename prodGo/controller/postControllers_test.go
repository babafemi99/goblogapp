package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"prodGo/cache"
	"prodGo/entity"
	"prodGo/repo"
	"prodGo/service"
	"strconv"

	// "strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	reposit        repo.PostRepository = repo.NewSqliteRepository()
	postSrv        service.PostService = service.NewPostService(reposit)
	PostCacheSrv cache.PostCache = cache.NewRedisCache("localhost:6379",0,10)
	postController PostController      = NewPostController(postSrv, PostCacheSrv)
)

const (
	ID = int64(123)
	TITLE string = "Title 1"
	Text string  = "Text 1"
)

func TestAddpost(t *testing.T) {
	// create a new http post request
	jsonData := []byte(
		`{
	"title": "Title 1", 
	"text":"Text 1"
	}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonData))
	// asing addpost as handler function for that request (cotroller addPost fn)
	handler := http.HandlerFunc(postController.AddPost)
	//record http response with http test library
	response := httptest.NewRecorder()
	//Dispatch the http request
	handler.ServeHTTP(response, req)
	// Add assertions on the http status code and response
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong error code : value is %v expected is %v", status, http.StatusOK)
	}
	// Decode the http response
	var post entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&post)
	// Assert HTTp response

	assert.NotNil(t, post.Id)
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, Text, post.Text)
	// clean up database
	//cod e to clean up database wold be here (remove all db)
}

// func cleanUp(post) {

// }

func TestGetPost(t *testing.T) {
	// insert new post 
	setup()
	// create a get http request
	req, _ := http.NewRequest("GET", "/posts", nil)
	// add post handler function for the request 
	handler := http.HandlerFunc(postController.GetPost)
	response  := httptest.NewRecorder()
	handler.ServeHTTP(response,req)
	status := response.Code
	if status != http.StatusOK{
		t.Errorf("Handler returned a wrong error code : value is %v expected is %v", status, http.StatusOK)
	}
	var posts []entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&posts)
	// Assert HTTP response
	assert.NotNil(t, posts[0].Id)
	assert.Equal(t, TITLE, posts[0].Title)
	assert.Equal(t, Text, posts[0].Text)
}

func setup(){
	var post entity.Post = entity.Post{
		Id: ID,
		Title: TITLE,
		Text: Text,
	}
	reposit.Save(&post)
}

func TestGetPostById (t *testing.T) {
	// insert new post 
	setup()
	// create a get http request
	req, _ := http.NewRequest("GET", "/posts/"+strconv.FormatInt(ID, 10), nil)
	// add post handler function for the request 
	handler := http.HandlerFunc(postController.GetPostById)
	response  := httptest.NewRecorder()
	handler.ServeHTTP(response,req)
	status := response.Code
	if status != http.StatusOK{
		t.Errorf("Handler returned a wrong error code : value is %v expected is %v", status, http.StatusOK)
	}
	var post entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&post)
	// Assert HTTP response
	assert.NotNil(t, ID, post.Id)
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, Text, post.Text)
}
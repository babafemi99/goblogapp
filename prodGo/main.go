package main

import (
	"fmt"
	"net/http"
	"os"

	// "prodGo/cache"
	"prodGo/cache"
	"prodGo/controller"
	"prodGo/repo"
	"prodGo/service"
	"prodGo/web"

	"github.com/joho/godotenv"
)

var (
	reposit repo.PostRepository = repo.NewDyanamoDBRepository()
)
var pService service.PostService = service.NewPostService(reposit)

var postCache cache.PostCache = cache.NewRedisCache("localhost:6379",1,10)
var webRouter web.RouterInterface = web.NewMuxRouter()

var postController controller.PostController = controller.NewPostController(pService, postCache)

func main() {
	godotenv.Load("prod.env")
	port := os.Getenv("PORT")
	webRouter.GET("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Up and running")
	})
	webRouter.GET("/post", postController.GetPost)
	webRouter.POST("/post", postController.AddPost)
	webRouter.GET("/post/{id}", postController.GetPostById)
	webRouter.SERVE(port)

}

package cache

import "prodGo/entity"

type PostCache interface {
	Set(key string, val *entity.Post)
	Get(key string) *entity.Post
}

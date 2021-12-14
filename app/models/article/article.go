package article

import (
	"goblog/pkg/route"
	"strconv"
)

// Article 文章模型
type Article struct {
	ID uint64
	Title string
	Body string
}

func (article Article) Link() string  {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(article.ID,10))
}

func GetStringID()  {
	
}

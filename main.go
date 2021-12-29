package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var router *mux.Router
var db *sql.DB

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们1,。</p>")
}

type Article struct {
	Title, Body string
	ID          int64
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//解析错误，错误处理
		fmt.Fprint(w, "请提交正确的数据！")
		return
	}
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(w, "title 的值为: %v <br>", title)
	fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
	fmt.Fprintf(w, "body 的值为: %v <br>", body)
	fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))
	fmt.Fprintf(w, "直接的值:%v<br>", r.PostFormValue("title"))
	fmt.Fprintf(w, "直接的值:%v<br>", r.PostFormValue("body"))
	//
	//fmt.Fprintf(w, "POST PostForm: %v<br>", r.PostForm)
	//fmt.Fprintf(w,"Form:%v<br>", r.Form)
	//fmt.Fprintf(w, "title的值: %v", title)
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取URL参数
	vars := mux.Vars(r)
	id := vars["id"]

	//2.读取对应的文字数据
	article := Article{}

	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	//3.出来错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1未找到数据
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404")
		} else {
			//3.2数据库错误
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500")
		}
	} else {
		//4.读取成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取参数
	id := getRouteVariable("id", r)

	//2.获取对应的文章数据
	article, err := getArticleById(id)
	fmt.Println(article)
	//3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
	} else {
		//4.未出现错误
		//4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			// 4.2 表单通过验证
			query := "update articles set title = ?,body = ? where id =?"
			res, err := db.Exec(query, title, body, id)

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			//成功更新，跳转到文章详情页
			if n, _ := res.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改")
			}
		} else {
			//4.3 验证未通过

			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			err = tmpl.Execute(w, data)
			logger.LogError(err)
		}
	}

}

func getArticleById(id string) (Article, error) {
	article := Article{}

	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

//中间件，设置标头
func forceHtmlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 1,设置标头
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		//2,继续处理
		next.ServeHTTP(writer, request)
	})
}

//中间件，
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
		}
		next.ServeHTTP(writer, request)
	})
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

}

func saveArticleToDB(title string, body string) (int64, error) {
	//变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	//1.获取一个prepare 声明语句
	stmt, err = db.Prepare("insert into articles (title,body) values (?,?)")
	//错误检查
	if err != nil {
		return 0, err
	}

	//2.在此函数运行结束后关闭此语句,防止占用SQL连接
	defer stmt.Close()

	//3.执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	//4.插入成功的话，会返回自增的ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	id := getRouteVariable("id", r)
	article,err := getArticleById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 未找到文章")
		}else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}else {
		// 执行删除
		rowsAffected,err :=article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误 500")
		}else {
			//无错误
			if rowsAffected > 0 {
				//重定向到文章列表
				indexURL,_ := router.Get("articles.index").URL()
				http.Redirect(w,r,indexURL.String(),http.StatusFound)
			}else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "未找到文章404")
			}
		}
	}
}

func (a Article) Delete()(rowsAffected int64,err error)  {
	rs,err := db.Exec("delete from articles where id =" + strconv.FormatInt(a.ID,10))
	if err != nil{
		return 0, err
	}

	//删除成功，跳转到文章详情页
	if n,_ := rs.RowsAffected();n > 0 {
		return n,nil
	}

	return 0, nil
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	//中间件
	router.Use(forceHtmlMiddleware)

	//通过命名路由获取URL
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("home:", homeUrl)

	//articleUrl, _ := router.Get("articles.show").URL("id", "23")
	//fmt.Println("articleUrl", articleUrl)

	http.ListenAndServe("127.0.0.1:3000", removeTrailingSlash(router))
}

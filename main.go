package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "mysql:3306",
		Net:                  "tcp",
		DBName:               "cesi",
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试链接，失败报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们1,。</p>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "唯有生活不可辜负！")
}

type Article struct {
	Title, Body string
	ID          int64
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取参数
	vars := mux.Vars(r)
	id := vars["id"]

	//2.读取文字数据
	article, err := getArticleById(id)

	//3.出现错误
	if err != nil {
		//3.1未找到数据
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}

	} else {
		//4.读取成功
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		err = tmpl.Execute(w, article)
		checkError(err)
	}
	fmt.Fprintf(w, "ID"+id)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>欢迎来到我的世界！</h1>")
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	//1.执行查询，返回结果集
	rows, err := db.Query("select * from articles")
	checkError(err)
	defer rows.Close()

	var articles []Article
	// 2.循环读取结果
	for rows.Next() {
		var article Article
		//2.1 扫描每一行的结果并赋值到一个article 对象中
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		//2.2将article 追加到articles的这个数组中
		articles = append(articles, article)
	}

	//2.3 检查遍历时是否发生错误
	err = rows.Err()
	checkError(err)

	// 3.加载模板
	tmpl,err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	//4.渲染模板，将所有的文章数据传输进去
	err = tmpl.Execute(w, articles)
	checkError(err)
}

func (a Article)Link() string  {
	showURL,err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID,10))
	if err != nil {
		checkError(err)
		return ""
	}
	return showURL.String()
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
			checkError(err)
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
		checkError(err)
		err = tmpl.Execute(w, data)
		checkError(err)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取参数
	id := getArticleVariable("id", r)

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
				checkError(err)
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
			checkError(err)
			err = tmpl.Execute(w, data)
			checkError(err)
		}
	}

}
func getArticleVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
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

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
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

func main() {
	initDB()
	createTables()

	//首页
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")

	//关于我们
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	//文章详情
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	//文章列表
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")

	//创建
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	//创建页面
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

	//404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//中间件
	router.Use(forceHtmlMiddleware)

	//通过命名路由获取URL
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("home:", homeUrl)

	//articleUrl, _ := router.Get("articles.show").URL("id", "23")
	//fmt.Println("articleUrl", articleUrl)

	http.ListenAndServe("127.0.0.1:3000", removeTrailingSlash(router))
}

package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

type User struct {
	ID    int    `gorm:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	db := gormConnect()
	defer db.Close()

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("views/templates/*.html")),
	}

	//CSSと画像のpath設定
	e.Static("/css", "./views/assets")
	e.Static("/image", "./views/assets")

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		var users []User
		// 全レコード取得
		db.Find(&users) //usersのアドレス情報を使って書き換えている
		return c.Render(http.StatusOK, "index", users)
	})

	e.GET("/:id", func(c echo.Context) error {
		var user User
		id := c.Param("id")
		db.Delete(&user, id)
		return c.Render(http.StatusOK, "delete", nil)
	})

	e.GET("/create", func(c echo.Context) error {
		return c.Render(http.StatusOK, "create", nil)
	})

	e.POST("/create/complete", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		user := User{Name: name, Email: email}
		db.Create(&user)
		return c.Render(http.StatusOK, "complete", nil)
	})

	e.Logger.Fatal(e.Start(":9000"))
}

func gormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"
	PASS := "root"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "go_crud"

	CONNECT := USER + ":" + "@" + PROTOCOL + "/" + DBNAME
	db, err := gorm.Open(DBMS, CONNECT)
	db.AutoMigrate(&User{})

	fmt.Println("SQL connecting...", USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME)
	if err != nil {
		panic(err.Error())
	}
	return db
}

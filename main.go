package main

import (
	// "echo.go"

	"context"
	"fmt"
	"log"
	"myportfolio/connection"
	"myportfolio/middleware"
	"net/http"
	"strconv"
	"text/template"
	"time"

	// "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	Description  string
	Technologies []string
	Image        string
	Duration     int64
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}

func main() {
	// create connection from database
	connection.DatabaseConnect()

	// fmt.Println(dataProject)

	// initialization echo
	e := echo.New()

	// server static from public directory
	e.Static("/public", "public")
	e.Static("/upload", "upload")

	// initialization to use session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/project", project)
	e.GET("/detail-project/:id", detail)
	e.GET("/delete/:id", deleteProject)
	e.GET("/form-register", formRegister)
	e.GET("/form-login", formLogin)
	e.GET("/logout", logout)
	e.GET("update/:id", updateFormProject)
	e.GET("/logout", logout)

	e.POST("/add-project", middleware.UploadFile(addProject))
	e.POST("/update-project", updateProject)
	// e.POST("/update-project", middleware.UploadFile(updateProject))
	e.POST("/register", register)
	e.POST("/login", login)

	e.Logger.Fatal(e.Start(":5000"))

}

func home(c echo.Context) error {

	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	// map(data type) => key and value
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects ORDER BY tb_projects.id DESC")
	// data, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, name, start_date, end_date, description, technologies, image, tb_user.username AS creator FROM tb_projects LEFT JOIN tb_user ON tb_projects.creator = tb_user.id ORDER BY tb_projects.id DESC")

	// create array for the result
	var result []Project

	for data.Next() {
		var each = Project{}

		// err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		firstDate := &each.StartDate
		secondDate := &each.EndDate
		difference := secondDate.Sub(*firstDate)

		d := int64(difference.Hours() / 24 / 30)

		each.Duration = d

		result = append(result, each)
	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	flash := map[string]interface{}{
		"Project":      result,
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
		"DataSession":  userData,
		"FlashID":      sess.Values["id"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}

func contact(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact_form.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func project(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/my_project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
		"DataSession":  userData,
	}

	return tmpl.Execute(c.Response(), flash)
}

func detail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/project_detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image  FROM tb_projects WHERE name = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	return tmpl.Execute(c.Response(), data)
}

func updateFormProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/update_project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func addProject(c echo.Context) error {
	c.Request().ParseForm()

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	flash := map[string]interface{}{
		"FlashID": sess.Values["id"],
	}

	flashID := flash["FlashID"].(int)

	id := flashID

	fmt.Println(id)

	_, err := connection.Conn.Exec(context.Background(), "SELECT id FROM public.tb_user WHERE id = $1;", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	nextJs := c.FormValue("nextJs")
	typescript := c.FormValue("typescript")
	description := c.FormValue("description")
	image := c.Get("dataFile").(string)
	creator := flashID

	var technologies []string

	if nodeJs == "nodejs" {
		technologies = append(technologies, "nodejs")
	} else {
		technologies = append(technologies, "")
	}
	if reactJs == "reactjs" {
		technologies = append(technologies, "reactjs")
	} else {
		technologies = append(technologies, "")
	}
	if nextJs == "nextjs" {
		technologies = append(technologies, "nextjs")
	} else {
		technologies = append(technologies, "")
	}
	if typescript == "typescript" {
		technologies = append(technologies, "typescript")
	} else {
		technologies = append(technologies, "")
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(name, start_date, end_date, description, technologies, image, creator) VALUES($1, $2, $3, $4, $5, $6, $7)", projectName, startDate, endDate, description, technologies, image, creator)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context) error {
	c.Request().ParseForm()
	id, _ := strconv.Atoi(c.Param("id"))

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	nextJs := c.FormValue("nextJs")
	typescript := c.FormValue("typescript")
	description := c.FormValue("description")
	// image := c.Get("dataFile").(string)
	image := "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mnx8dGVjaG5vbG9neXxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60"

	// var ProjectDetail = Project{}

	var technologies []string

	if nodeJs == "nodejs" {
		technologies = append(technologies, "nodejs")
	} else {
		technologies = append(technologies, "")
	}
	if reactJs == "reactjs" {
		technologies = append(technologies, "reactjs")
	} else {
		technologies = append(technologies, "")
	}
	if nextJs == "nextjs" {
		technologies = append(technologies, "nextjs")
	} else {
		technologies = append(technologies, "")
	}
	if typescript == "typescript" {
		technologies = append(technologies, "typescript")
	} else {
		technologies = append(technologies, "")
	}

	// id := ProjectDetail.Id

	// err := connection.Conn.QueryRow(context.Background(), "SELECT id FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id)

	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	// }

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7", projectName, startDate, endDate, description, technologies, image, id)

	fmt.Println(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	// _, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7", projectName, startDate, endDate, description, technologies, image, id)

	// fmt.Println(id)

	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	// }

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func formRegister(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")

	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return redirectWithMessage(c, "Email Salah !", false, "/form-login")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 jam
	sess.Values["message"] = "Login success"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true // access login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// generate password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (username, email, password) VALUES($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		redirectWithMessage(c, "Register failed, please try again :)", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success", true, "/form-login")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Values["isLogin"] = false
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}

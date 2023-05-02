package main

import (
	// "echo.go"

	"context"
	"fmt"
	"myportfolio/connection"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
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
	Char         string
}

var dataProject = []Project{
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		//StartDate:   "13/04/2020",
		//EndDate:     "13/04/2023",
	},
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		// StartDate:   "13/04/2020",
		// EndDate:     "13/04/2023",
	},
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		// StartDate:   "13/04/2020",
		// EndDate:     "13/04/2023",
	},
}

func main() {
	// create connection from database
	connection.DatabaseConnect()

	// fmt.Println(dataProject)

	// initialization echo
	e := echo.New()

	// server static from public directory
	e.Static("/public", "public")

	// routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/project", project)
	e.GET("/detail-project/:id", detail)
	e.POST("/add-project", addProject)
	e.GET("/delete/:id", deleteProject)
	e.GET("login", login)
	e.GET("register", register)

	e.Logger.Fatal(e.Start(":5000"))

}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		// fmt.Println("Page not found")
		// return c.String(http.StatusOK, "Hello World!!")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	// map(data type) => key and value
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects")

	// create array for the result
	var result []Project

	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		// each.Author = ""

		firstDate := &each.StartDate
		secondDate := &each.EndDate
		difference := secondDate.Sub(*firstDate)

		// fmt.Printf("Months: %d\n", int64(difference.Hours()/24/30/12))
		// fmt.Printf("Months: %d\n", int64(difference.Hours()/24/30))
		// fmt.Printf("Months: %d\n", int64(difference.Hours()/24))
		// fmt.Printf("Months: %d\n", int64(difference.Hours()))

		d := int64(difference.Hours() / 24 / 30)

		each.Duration = d

		// t := each.Technologies
		// for _, char := range t {
		// 	switch char {
		// 	case "reactjs":
		// 		// p := `<img src="/public/img/nodejs.png>`
		// 		// return p
		// 		fmt.Println(`<img src="/public/img/reactjs.png>`)
		// 	case "nodejs":
		// 		fmt.Println(`<img src="/public/img/nodejs.png>`)
		// 	case "nextjs":
		// 		fmt.Println(`<img src="/public/img/nextJs.png>`)
		// 	case "typescript":
		// 		fmt.Println(`<img src="/public/img/ts.png>`)
		// 	default:
		// 		fmt.Println("")
		// 	}
		// }

		// for _, char := range t {
		// 	if char == "reactjs" {
		// 		return ""
		// 	}
		// }

		result = append(result, each)
	}

	// fmt.Printf("%v", result)

	projects := map[string]interface{}{
		"Project": result,
	}

	// projects := map[string]interface{}{
	// 	"Project": dataProject,
	// }

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact_form.html")

	if err != nil {
		// fmt.Println("Page not found")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func project(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/my_project.html")

	if err != nil {
		// fmt.Println("Page not found")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func detail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/project_detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"Id":          id,
		"ProjectName": "Aplikasi sederhana yang dibangun oleh golang",
		"Description": "saya merasakn apa yang anda rasakan sekarang",
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	description := c.FormValue("description")

	var addProject = Project{
		ProjectName: projectName,
		Description: description,
	}
	// fmt.Println(projectName)
	// fmt.Println(description)

	fmt.Println(addProject)
	dataProject = append(dataProject, addProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func login(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

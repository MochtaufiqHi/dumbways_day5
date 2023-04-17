package main

import (
	// "echo.go"

	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Project struct {
	ProjectName string
	Description string
	StartDate   string
	EndDate     string
}

var dataProject = []Project{
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		StartDate:   "13/04/2020",
		EndDate:     "13/04/2023",
	},
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		StartDate:   "13/04/2020",
		EndDate:     "13/04/2023",
	},
	{
		ProjectName: "Aplikasi Sederhana",
		Description: "Aplikasi sederhana yang dibangun menggunakan Javascript dan Boostraps",
		StartDate:   "13/04/2020",
		EndDate:     "13/04/2023",
	},
}

func main() {
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

	e.Logger.Fatal(e.Start(":5000"))
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		// fmt.Println("Page not found")
		// return c.String(http.StatusOK, "Hello World!!")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	projects := map[string]interface{}{
		"Project": dataProject,
	}

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

package main

import (
	// "echo.go"

	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

func main() {
	// initialization echo
	e := echo.New()

	// server static from public directory
	e.Static("/public", "public")

	// routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/project", project)
	e.GET("/project-detail/:id", projectDetail)
	e.POST("/add-project", addProject)

	e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		// fmt.Println("Page not found")
		// return c.String(http.StatusOK, "Hello World!!")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
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

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // string to int

	tmpl, err := template.ParseFiles("views/project_detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"Id":      id,
		"Tittle":  "Simple Calculator",
		"Author":  "Mochamamd Taufiq Hidayat",
		"Content": "Pengembangan aplikasi yang menngunakan bahasa GO yang sekarang sangat gencar di perbincangkan.",
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	description := c.FormValue("description")

	fmt.Println(projectName)
	fmt.Println(description)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

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
	e.GET("/login", login)
	e.GET("/register", register)

	e.GET("/update/:id", updateProject)

	e.Logger.Fatal(e.Start(":5000"))

}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
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

		firstDate := &each.StartDate
		secondDate := &each.EndDate
		difference := secondDate.Sub(*firstDate)

		d := int64(difference.Hours() / 24 / 30)

		each.Duration = d

		result = append(result, each)
	}

	projects := map[string]interface{}{
		"Project": result,
	}

	return tmpl.Execute(c.Response(), projects)
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

	return tmpl.Execute(c.Response(), nil)
}

func detail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/project_detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image  FROM tb_projects WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	c.Request().ParseForm()

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	nextJs := c.FormValue("nextJs")
	typescript := c.FormValue("typescript")
	description := c.FormValue("description")
	image := "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8M3x8dGVjaHxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60"

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

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(name, start_date, end_date, description, technologies, image) VALUES($1, $2, $3, $4, $5, $6)", projectName, startDate, endDate, description, technologies, image)

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
	id, _ := strconv.Atoi(c.Param("id"))

	c.Request().ParseForm()

	projectName := c.FormValue("projectName")
	// startDate := c.FormValue("start-date")
	// endDate := c.FormValue("end-date")
	nodeJs := c.FormValue("nodeJs")
	reactJs := c.FormValue("reactJs")
	nextJs := c.FormValue("nextJs")
	typescript := c.FormValue("typescript")
	description := c.FormValue("description")
	image := "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8M3x8dGVjaHxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60"

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

	// _, err := connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6", projectName, startDate, endDate, description,
	// technologies, image)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$1, description=$2, technologies=$3, image=$4 WHERE id=$5", projectName, description, technologies, image, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/project")
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

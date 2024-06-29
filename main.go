package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)                 // Número máximo de conexões abertas ao mesmo tempo
	db.SetMaxIdleConns(25)                 // Número máximo de conexões inativas no pool
	db.SetConnMaxLifetime(5 * time.Minute) // Tempo máximo de vida de uma conexão
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok!")
	})
	e.POST("/projects", func(c echo.Context) error {
		type Project struct {
			Title string `json:"title"`
		}
		var project Project
		err := c.Bind(&project)
		if err != nil {
			return err
		}
		if project.Title == "" {
			return c.JSON(http.StatusBadRequest, "title is required")
		}
		id := uuid.New()
		res, err := db.Exec("INSERT INTO project (id, title) VALUES (?, ?)", id.String(), project.Title)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			return err
		}
		c.Response().Header().Set("Location", "/projects/"+id.String())
		return c.JSON(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

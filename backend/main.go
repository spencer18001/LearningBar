package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

type LearningSession struct {
	Attempt   int  `json:"attempt"`
	Completed bool `json:"completed"`
}

type LearningItem struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Outline  string            `json:"outline"`
	Sessions []LearningSession `json:"sessions,omitempty"`
}

var (
	items   = []LearningItem{}
	itemMux sync.Mutex
)

func getAllItems(c echo.Context) error {
	itemMux.Lock()
	defer itemMux.Unlock()
	return c.JSON(http.StatusOK, items)
}

func createItem(c echo.Context) error {
	item := new(LearningItem)
	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	itemMux.Lock()
	defer itemMux.Unlock()
	item.ID = len(items) + 1
	items = append(items, *item)

	return c.JSON(http.StatusCreated, item)
}

func createSession(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	session := new(LearningSession)
	if err := c.Bind(session); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	itemMux.Lock()
	defer itemMux.Unlock()

	for i := range items {
		if items[i].ID == id {
			session.Attempt = len(items[i].Sessions) + 1
			items[i].Sessions = append(items[i].Sessions, *session)
			return c.JSON(http.StatusCreated, items[i])
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
}

func main() {
	e := echo.New()

	e.GET("/api/items", getAllItems)
	e.POST("/api/items", createItem)
	e.POST("/api/items/:id/sessions", createSession)

	e.Logger.Fatal(e.Start(":8080"))
}

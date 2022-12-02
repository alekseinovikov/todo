package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"todo/service"
)

func NewTodoController(e *echo.Echo, s service.TodoService) TodoController {
	return &todoController{e: e, s: s}
}

type todoController struct {
	e *echo.Echo
	s service.TodoService
}

func (t *todoController) RegisterRoutes() {
	t.initStatic()
	t.initCrud()
}

func (t *todoController) initCrud() {
	group := t.e.Group("/api/todos")

	group.GET("/:id", t.getById)
	group.POST("", t.create)
	group.PUT("/:id", t.update)
	group.POST("/markDone/:id", t.markDone)
	group.POST("/markUndone/:id", t.markUndone)
}
func (t *todoController) initStatic() {
	t.e.Static("/", "static")
	t.e.File("/", "static/index.html")
}

func (t *todoController) getById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	todo, err := t.s.FindById(uint32(id))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if todo.IsAbsent() {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo.MustGet())
}

func (t *todoController) create(c echo.Context) error {
	var ct service.CreateTodo
	err := c.Bind(&ct)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	todo, err := t.s.Save(ct)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, todo)
}

func (t *todoController) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var ut service.UpdateTodo
	err = c.Bind(&ut)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	todo, err := t.s.Update(uint32(id), ut)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, todo)
}

func (t *todoController) markDone(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = t.s.MarkDone(uint32(id))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (t *todoController) markUndone(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = t.s.MarkUndone(uint32(id))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mircearem/storer/store"
)

type Server struct {
	e     *echo.Echo
	db    *store.Store
	errch chan error
}

func NewServer(db *store.Store) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(http.StatusInternalServerError, store.Map{"error": err.Error()})
	}
	return &Server{
		e:     e,
		db:    db,
		errch: make(chan error),
	}
}

func (s *Server) Run() error {
	// Setup API endpoints
	s.e.POST("/api/:collname", s.handlePostInsert)
	s.e.GET("/api/:collname", s.handleGetQuery)
	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		if err := s.e.Start(port); err != nil {
			s.errch <- err
		}
	}()
	return <-s.errch
}

func (s *Server) handlePostInsert(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
		id       uint64
	)
	var data store.Map
	if err = json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return err
	}
	for k, v := range data {
		if id, err = s.db.Collection(collname).Put([]byte(k), []byte(v)); err != nil {
			return err
		}
	}
	return c.JSON(http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleGetQuery(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
	)

	jmap := make(map[string]string)
	if err = json.NewDecoder(c.Request().Body).Decode(&jmap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if _, ok := jmap["key"]; !ok {
		msg := "\"key\" tag not present"
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}
	rec, err := s.db.Collection(collname).Get([]byte(jmap["key"]))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"value": string(rec)})
}

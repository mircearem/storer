package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mircearem/storer/log"
	"github.com/mircearem/storer/store"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	e     *echo.Echo
	db    *store.Store
	errch chan error
}

func NewApiServer(db *store.Store) *ApiServer {
	log := logrus.New()
	apiLogger := NewApiLogger(log)

	e := echo.New()
	e.Logger = apiLogger
	e.Use(middleware.Logger())

	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(http.StatusInternalServerError, store.Map{"error": err.Error()})
	}
	return &ApiServer{
		e:     e,
		db:    db,
		errch: make(chan error),
	}
}

func (s *ApiServer) Run() error {
	s.e.GET("/store/:collname", s.handleGetQuery)
	s.e.POST("/store/:collname", s.handlePostQuery)
	s.e.DELETE("/store/:collname", s.handleDeleteQuery)

	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		if err := s.e.Start(port); err != nil {
			s.errch <- err
		}
	}()
	return <-s.errch
}

func (s *ApiServer) handleGetQuery(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
	)
	var jmap store.Map
	if err = json.NewDecoder(c.Request().Body).Decode(&jmap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	key, ok := jmap["key"]
	if !ok {
		msg := "\"key\" tag not present"
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}
	rec, err := s.db.Collection(collname).Get([]byte(key))
	if err != nil {
		switch err.(store.StoreError).Type() {
		case store.ERR_COL_NOT_FOUND:
			return c.JSON(http.StatusNoContent, map[string]string{"error": err.Error()})
		case store.ERR_GET_FAIL_UDEF:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		case store.ERR_GET_FAIL_NOTF:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusFound, map[string]string{"value": string(rec)})
}

func (s *ApiServer) handlePostQuery(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
		id       uint64
	)
	var message store.Message
	if err = json.NewDecoder(c.Request().Body).Decode(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	id, err = s.db.Collection(collname).Set([]byte(message.Key), []byte(message.Value), message.TTL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]uint64{"id": id})
}

func (s *ApiServer) handleDeleteQuery(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
	)
	var jmap store.Map
	if err = json.NewDecoder(c.Request().Body).Decode(&jmap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	key, ok := jmap["key"]
	if !ok {
		msg := "\"key\" tag not present"
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}
	err = s.db.Collection(collname).Delete([]byte(key))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

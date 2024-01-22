package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mircearem/storer/store"
	"github.com/sirupsen/logrus"
)

type ApiServer struct {
	e       *echo.Echo
	storage *store.StorageServer
	errch   chan error
}

func NewApiServer(storage *store.StorageServer) *ApiServer {
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
		e:       e,
		storage: storage,
		errch:   make(chan error),
	}
}

func (s *ApiServer) Run() error {
	s.e.POST("/store/:collname", s.handlePostInsert)
	s.e.GET("/store/:collname", s.handleGetQuery)

	// go s.storage.Run()

	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		if err := s.e.Start(port); err != nil {
			close(s.storage.Closech)
			s.errch <- err
		}
	}()
	return <-s.errch
}

func (s *ApiServer) handlePostInsert(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
		id       uint64
	)
	var data store.Map
	if err = json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	for k, v := range data {
		id, err = s.storage.DB.Collection(collname).Put([]byte(k), []byte(v))
		if err != nil {
			dberr := err.(store.StoreError)
			switch dberr.Type() {
			// err_put_fail_undefined
			case store.ERR_PUT_FAIL_UDEF:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": dberr.Error()})
			// err_put_fail_conflict
			case store.ERR_PUT_FAIL_CONF:
				return c.JSON(http.StatusConflict, map[string]string{"error": dberr.Error()})
			}
		}
	}
	return c.JSON(http.StatusCreated, map[string]uint64{"id": id})
}

func (s *ApiServer) handleGetQuery(c echo.Context) error {
	var (
		collname = c.Param("collname")
		err      error
	)
	// key := c.QueryParam("key")
	// if key == "" {
	// 	msg := "\"key\" query param not present"
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	// }
	var jmap store.Map
	if err = json.NewDecoder(c.Request().Body).Decode(&jmap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	key, ok := jmap["key"]
	if !ok {
		msg := "\"key\" tag not present"
		return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
	}
	rec, err := s.storage.DB.Collection(collname).Get([]byte(key))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"value": string(rec)})
}

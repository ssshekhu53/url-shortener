package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"url-shortener/models"
	"url-shortener/service"
)

type URL interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetAnalytics(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type url struct {
	service service.URL
}

func New(service service.URL) URL {
	return &url{service: service}
}

func (u *url) Create(c *gin.Context) {
	var req models.CreateRequest

	err := c.Bind(&req)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	err = req.Validate()
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	alias, err := u.service.Create(&req)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	resp := models.CreateResponse{ShortURL: alias}

	c.JSON(http.StatusCreated, resp)
}

func (u *url) Get(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, "missing alias")

		return
	}

	longURL, err := u.service.Get(alias)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	c.Redirect(http.StatusTemporaryRedirect, longURL)
}

func (u *url) GetAnalytics(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, "missing alias")

		return
	}

	analytics, err := u.service.GetAnalytics(alias)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (u *url) Update(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, "missing alias")

		return
	}

	var req models.UpdateRequest

	err := c.Bind(&req)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	err = req.Validate()
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	err = u.service.Update(&req, alias)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	c.JSON(http.StatusOK, "Successfully updated")
}

func (u *url) Delete(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, "missing alias")

		return
	}

	err := u.service.Delete(alias)
	if err != nil {
		u.errorRespons(c, err)

		return
	}

	c.JSON(http.StatusOK, "Successfully deleted")
}

func (u *url) errorRespons(c *gin.Context, err error) {
	var statusCode int
	errorMessage := err.Error()

	if strings.Contains(errorMessage, "alias does not exist or has expired") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(errorMessage, "invalid request") || strings.Contains(errorMessage, "missing") {
		statusCode = http.StatusBadRequest
	} else if strings.Contains(errorMessage, "already exists") {
		statusCode = http.StatusConflict
	} else {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, errorMessage)
}

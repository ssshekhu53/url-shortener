package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ssshekhu53/url-shortener/handler"
	"github.com/ssshekhu53/url-shortener/service"
	"github.com/ssshekhu53/url-shortener/store"
	"github.com/ssshekhu53/url-shortener/utils"
)

func main() {
	app := gin.Default()

	urlStore := store.New()
	urlService := service.New(urlStore)
	urlHandler := handler.New(urlService)

	aliasRemoverObj := &aliasRemover{service: urlService}

	go aliasRemoverObj.Remove()

	app.POST("/shorten", urlHandler.Create)
	app.GET("/:alias", urlHandler.Get)
	app.GET("/analytics/:alias", urlHandler.GetAnalytics)
	app.PUT("/update/:alias", urlHandler.Update)
	app.DELETE("/delete/:alias", urlHandler.Delete)

	app.Run(":8000")
}

type aliasRemover struct {
	service service.URL
}

func (a *aliasRemover) Remove() {
	for {
		analytics := a.service.GetAllAnalytics()

		for i := range analytics {
			hasExpired, _ := utils.HasExpired(&analytics[i])

			if hasExpired {
				a.service.Delete(analytics[i].Alias)
			}
		}
	}
}

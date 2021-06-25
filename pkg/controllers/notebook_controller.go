package controllers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elastic-jupyter-dashboard-server/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetNotebooks(ctx *gin.Context) {
	notebookModel := models.Notebook{}

	if data, err := notebookModel.GetNotebooks(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
	} else {
		if len(data) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"data": data,
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"data": data,
			})
		}
	}
}

func DeleteNotebook(ctx *gin.Context) {
	name := ctx.Query("name")
	namespace := ctx.Query("namespace")
	notebookModel := models.Notebook{}

	if err := notebookModel.DeleteNotebook(name, namespace); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
	} else {
		ctx.Status(http.StatusOK)
	}
}

func CreateNotebook(ctx *gin.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	log.Println(jsonData)

	notebookModel := models.Notebook{}

	if err := ctx.ShouldBind(&notebookModel); err == nil {
		log.Println(notebookModel.Name)
		log.Println(notebookModel.Namespace)
		log.Println("bind")
	} else {
		log.Println("error")
	}
}

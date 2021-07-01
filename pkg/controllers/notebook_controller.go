package controllers

import (
	"net/http"

	"github.com/elastic-jupyter-dashboard-server/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	notebookModel := models.Notebook{}
	var err error
	if err = ctx.ShouldBindWith(&notebookModel, binding.JSON); err == nil {
		if err = notebookModel.CreateNotebook(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		} else {
			ctx.Status(http.StatusCreated)
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
}

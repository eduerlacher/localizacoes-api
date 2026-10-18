package controller

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"localizacoes-api/internal/model"
	"localizacoes-api/internal/view"
)

// UnidadeModel descreve as operações de dados usadas pelos controllers.
type UnidadeModel interface {
	ListUPAs(ctx context.Context) ([]model.UPA, error)
	ListUBSs(ctx context.Context) ([]model.UBS, error)
}

// UnidadeController trata as requisições HTTP relacionadas às unidades.
type UnidadeController struct {
	model UnidadeModel
}

func NewUnidadeController(m UnidadeModel) *UnidadeController {
	return &UnidadeController{model: m}
}

func (c *UnidadeController) ListUPAs(ctx *gin.Context) {
	upas, err := c.model.ListUPAs(ctx.Request.Context())
	if err != nil {
		c.handleDataError(ctx, err, "UPA", "Failed to read UPA", "Failed to read UPA rows")
		return
	}

	ctx.JSON(http.StatusOK, view.UPAs(upas))
}

func (c *UnidadeController) ListUBSs(ctx *gin.Context) {
	ubss, err := c.model.ListUBSs(ctx.Request.Context())
	if err != nil {
		c.handleDataError(ctx, err, "UBS", "Failed to read UBS", "Failed to read UBS records")
		return
	}

	ctx.JSON(http.StatusOK, view.UBSs(ubss))
}

func (c *UnidadeController) handleDataError(ctx *gin.Context, err error, resource, scanMessage, rowsMessage string) {
	var dataErr *model.DataError
	if !errors.As(err, &dataErr) {
		log.Printf("Erro ao acessar dados de %s: %v", resource, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}

	switch dataErr.Stage {
	case model.QueryStage:
		log.Printf("ERRO SQL %s: %v", resource, dataErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database query failed",
			"details": dataErr.Error(),
		})
	case model.ScanStage:
		log.Printf("ERRO SCAN %s: %v", resource, dataErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": scanMessage})
	case model.RowsStage:
		if resource == "UPA" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": rowsMessage})
			return
		}
		log.Printf("%s rows error: %v", resource, dataErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": rowsMessage})
	default:
		log.Printf("Erro ao acessar dados de %s: %v", resource, dataErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": dataErr.Error()})
	}
}

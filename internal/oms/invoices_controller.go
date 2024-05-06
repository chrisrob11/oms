package oms

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chrisrob11/oms/internal/oms/db"
	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/gin-gonic/gin"
)

type invoicesController struct {
	dbQueries *db.Queries
	logger    *slog.Logger
}

func newInvoicesController(logger *slog.Logger, engine *gin.Engine, dbQueries *db.Queries) *invoicesController {
	controller := &invoicesController{dbQueries: dbQueries, logger: logger}
	engine.POST("/invoices", controller.create)
	engine.GET("/invoices", controller.list)
	engine.GET("/invoices/:id", controller.get)
	engine.DELETE("/invoices/:id", controller.delete)
	engine.POST("/invoices/:id/adjust", controller.adjust)

	return controller
}

func (s *invoicesController) create(c *gin.Context) {
	var req models.Invoice
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceID, err := s.dbQueries.CreateInvoice(c.Request.Context(), req.ToCreateInvoiceParams())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoiceID)
}

func (s *invoicesController) get(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	invoice, err := s.dbQueries.GetInvoice(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	invoiceModel, err := models.NewInvoiceFromDB(invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoiceModel)
}

func (s *invoicesController) list(c *gin.Context) {
	params := db.ListInvoicesParams{
		Limit: 100,
	}

	limitInt, hasError := extractLimit(c)
	if hasError {
		return
	}

	if limitInt != nil {
		params.Limit = *limitInt
	}

	pageInfo, hasError := extractTokenFromQuery(c)
	if hasError {
		return
	}

	if pageInfo != nil {
		s.logger.Info("Token info", slog.Attr{Key: "starting_id", Value: slog.IntValue(pageInfo.StartID)})
		params.ID = int32(pageInfo.StartID)
		params.Limit = int32(pageInfo.Size)
	}

	invoices, err := s.dbQueries.ListInvoices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	numItems := len(invoices)
	invoicesResp := &models.List[models.Invoice]{}
	invoicesResp.Items = make([]*models.Invoice, numItems)

	for i := 0; i < numItems; i++ {
		v, err := models.NewInvoiceFromDB(invoices[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		invoicesResp.Items[i] = v
	}

	if numItems >= int(params.Limit) {
		token := EncodeToken(PaginationToken{StartID: int(invoices[numItems-1].ID), Size: int(params.Limit)})
		invoicesResp.NextPageToken = token
	}

	c.JSON(http.StatusOK, invoicesResp)
}

func (s *invoicesController) adjust(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var adjustmentAmount float64
	if err := c.BindJSON(&adjustmentAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adjustedAmountStr := strconv.FormatFloat(adjustmentAmount, 'f', 15, 64)
	params := db.AdjustInvoiceParams{ID: id, TotalAdjustments: sql.NullString{Valid: true, String: adjustedAmountStr}}

	if err := s.dbQueries.AdjustInvoice(c.Request.Context(), params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, int(id))
}

func (s *invoicesController) delete(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = s.dbQueries.DeleteInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

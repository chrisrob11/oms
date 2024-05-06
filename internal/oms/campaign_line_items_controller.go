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

type campaignLineItemsController struct {
	dbQueries *db.Queries
	logger    *slog.Logger
}

func newCampaignLineItemsController(logger *slog.Logger, engine *gin.Engine,
	dbQueries *db.Queries) *campaignLineItemsController {
	controller := &campaignLineItemsController{dbQueries: dbQueries, logger: logger}
	engine.POST("/campaignLineItems", controller.create)
	engine.GET("/campaignLineItems", controller.list)
	engine.GET("/campaignLineItems/:id", controller.get)
	engine.PUT("/campaignLineItems/:id", controller.update)
	engine.DELETE("/campaignLineItems/:id", controller.delete)

	return controller
}

func (s *campaignLineItemsController) create(c *gin.Context) {
	var modelReq models.CampaignLineItem
	if err := c.BindJSON(&modelReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var campaignLineID int32

	var err error

	if modelReq.ID > 0 {
		req := modelReq.ToCreateCampaignLineItemWithID()

		campaignLineID, err = s.dbQueries.CreateCampaignLineWithID(c.Request.Context(), req)
		if err == nil {
			resetErr := s.dbQueries.ResetCampaignLineItemID(c.Request.Context())
			if resetErr != nil {
				s.logger.Error("resetting the serial id for campaign line item failed")
			}
		}
	} else {
		req := modelReq.ToCreateCampaignLineItem()
		campaignLineID, err = s.dbQueries.CreateCampaignLine(c.Request.Context(), req)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, campaignLineID)
}

func (s *campaignLineItemsController) get(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	campaignLine, err := s.dbQueries.GetCampaignLine(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clm, err := models.NewCampaignLineItemFromDB(&campaignLine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clm)
}

func (s *campaignLineItemsController) list(c *gin.Context) {
	params := db.ListCampaignLineItemsParams{
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

	campaignLineItems, err := s.dbQueries.ListCampaignLineItems(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	numItems := len(campaignLineItems)
	campaignLineItemsResp := &models.List[models.CampaignLineItem]{}
	campaignLineItemsResp.Items = make([]*models.CampaignLineItem, numItems)

	for i := 0; i < numItems; i++ {
		v, err := models.NewCampaignLineItemFromDB(&campaignLineItems[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		campaignLineItemsResp.Items[i] = v
	}

	if numItems >= int(params.Limit) {
		token := EncodeToken(PaginationToken{StartID: int(campaignLineItems[numItems-1].ID), Size: int(params.Limit)})
		campaignLineItemsResp.NextPageToken = token
	}

	c.JSON(http.StatusOK, campaignLineItemsResp)
}

func (s *campaignLineItemsController) update(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req models.CampaignLineItem
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	col := req.ToCreateCampaignLineItem()
	update := db.UpdateCampaignLineParams{
		ID:          id,
		Name:        req.Name,
		Booked:      col.Booked,
		Actual:      col.Actual,
		Adjustments: col.Adjustments,
		StartedAt:   col.StartedAt,
		EndedAt:     col.EndedAt,
	}

	if err := s.dbQueries.UpdateCampaignLine(c.Request.Context(), update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *campaignLineItemsController) delete(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = s.dbQueries.DeleteCampaignLine(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func toInt32(v string) (int32, error) {
	id64, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(id64), nil
}

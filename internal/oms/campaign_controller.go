package oms

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/chrisrob11/oms/internal/oms/db"
	"github.com/chrisrob11/oms/internal/oms/models"
	"github.com/gin-gonic/gin"
)

type campaignsController struct {
	dbQueries *db.Queries
	logger    *slog.Logger
}

func newCampaignsController(logger *slog.Logger, engine *gin.Engine, dbQueries *db.Queries) *campaignsController {
	controller := &campaignsController{dbQueries: dbQueries, logger: logger}
	engine.POST("/campaigns", controller.create)
	engine.GET("/campaigns/:id", controller.get)
	engine.GET("/campaigns", controller.list)
	engine.PUT("/campaigns/:id", controller.update)
	engine.POST("/campaigns/:id/generateInvoice", controller.generateInvoice)
	engine.DELETE("/campaigns/:id", controller.delete)

	return controller
}

func (s *campaignsController) create(c *gin.Context) {
	var req models.Campaign
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error

	var campaign db.OmsCampaign

	if req.ID > 0 {
		s.logger.Info("Creating campaign with id")

		params := req.ToCreateCampaignWithID()

		campaign, err = s.dbQueries.CreateCampaignWithID(c.Request.Context(), *params)
		if err == nil {
			resetErr := s.dbQueries.ResetCampaignID(c.Request.Context())
			if resetErr != nil {
				s.logger.Error("Cannot reset the serial id after manual insert")
			}
		}
	} else {
		s.logger.Info("Creating campaign without id")
		params := req.ToCreateCampaign()
		campaign, err = s.dbQueries.CreateCampaign(c.Request.Context(), *params)
	}

	if err != nil {
		s.logger.Error("Error creating campaign", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, models.NewCampaignFromDB(&campaign))
}

func (s *campaignsController) get(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	campaign, err := s.dbQueries.GetCampaign(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if err != nil {
		s.logger.Error("error occurred getting campaign", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, models.NewCampaignFromDB(&campaign))
}

func (s *campaignsController) list(c *gin.Context) {
	params := db.ListCampaignsParams{
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

	campaigns, err := s.dbQueries.ListCampaigns(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	numItems := len(campaigns)
	campaignsResp := &models.List[models.Campaign]{}
	campaignsResp.Items = make([]*models.Campaign, len(campaigns))

	for i := 0; i < len(campaigns); i++ {
		campaignsResp.Items[i] = models.NewCampaignFromDB(&campaigns[i])
	}

	if numItems >= int(params.Limit) {
		token := EncodeToken(PaginationToken{StartID: int(campaigns[numItems-1].ID), Size: int(params.Limit)})
		campaignsResp.NextPageToken = token
	}

	c.JSON(http.StatusOK, campaignsResp)
}

func (s *campaignsController) update(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req models.Campaign
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaignDB := req.ToCreateCampaign()

	// Have the problem where if some parameters are not specified
	// then empty will overwrite the value. Have another problem
	// where if empty is the value we don't know to set. There needs
	// to be more changes to fix these bugs.
	update := db.UpdateCampaignParams{
		Name:      campaignDB.Name,
		StartedAt: campaignDB.StartedAt,
		EndedAt:   campaignDB.EndedAt,
		Archiving: campaignDB.Archiving,
		ID:        id,
	}

	if err := s.dbQueries.UpdateCampaign(c.Request.Context(), update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *campaignsController) delete(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = s.dbQueries.DeleteCampaign(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *campaignsController) generateInvoice(c *gin.Context) {
	id, err := toInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	s.logger.Info("Building Campaign Invoice")
	invoice, err := s.dbQueries.BuildCampaignInvoice(c.Request.Context(), db.BuildCampaignInvoiceParams{
		CampaignID: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Creating the Invoice")

	invoice.IssuedAt = sql.NullTime{Valid: true, Time: time.Now().UTC()}

	createdID, err := s.dbQueries.CreateInvoice(c.Request.Context(), invoice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	c.JSON(http.StatusOK, createdID)
}

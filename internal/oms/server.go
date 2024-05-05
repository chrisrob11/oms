package oms

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/chrisrob11/oms/internal/oms/db"
)

var DB *sql.DB

type Server struct {
	db                      *db.Queries
	logger                  *slog.Logger
	engine                  *gin.Engine
	campaignsController     *campaignsController
	campaignlinesController *campaignLineItemsController
	invoicesController      *invoicesController
}

func NewServer() (*Server, error) {
	r := gin.Default()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	db, err := db.NewDBContext()
	if err != nil {
		return nil, err
	}

	campaignController := newCampaignsController(logger, r, db)
	campaignLineItemsController := newCampaignLineItemsController(logger, r, db)
	invoices := newInvoicesController(logger, r, db)

	return &Server{engine: r, db: db, campaignsController: campaignController,
		invoicesController: invoices, campaignlinesController: campaignLineItemsController,
		logger: logger,
	}, nil
}

func (s *Server) Run(ctx context.Context) (err error) {
	s.logger.Info("Server Starting")
	return s.engine.Run() // listen and serve on 0.0.0.0:8080
}

package handler

import (
	"database/sql"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/service"
)

type Handler struct {
	db                *sql.DB
	queries           *db.Queries
	fileService       service.FileService
	documentProcessor *service.DocumentProcessor
	searchService     *service.SearchService
	chatService       *service.ChatService
	analysisService   *service.AnalysisService
	sourceService     *service.SourceService
}

func NewHandler(
	database *sql.DB,
	fileService service.FileService,
	documentProcessor *service.DocumentProcessor,
	searchService *service.SearchService,
	chatService *service.ChatService,
	analysisService *service.AnalysisService,
	sourceService *service.SourceService,
) *Handler {
	return &Handler{
		db:                database,
		queries:           db.New(database),
		fileService:       fileService,
		documentProcessor: documentProcessor,
		searchService:     searchService,
		chatService:       chatService,
		analysisService:   analysisService,
		sourceService:     sourceService,
	}
}

// Compile-time check that Handler implements api.ServerInterface
var _ api.ServerInterface = (*Handler)(nil)

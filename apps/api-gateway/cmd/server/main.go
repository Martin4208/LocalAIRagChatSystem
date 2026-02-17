package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/config"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/handler"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/middleware"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/service"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/storage"
)

func main() {
	// Load configuration
	_ = godotenv.Load()
	cfg := config.Load()

	// log.Println("✅ UniDoc license activated")

	// --- Database Setup ---
	connStr := "host=" + cfg.Database.Host +
		" port=" + cfg.Database.Port +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.DBName +
		" sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	log.Println("✅ Database connected successfully")

	queries := db.New(database)

	// --- MinIO Setup ---
	minioEndpoint := cfg.MinIO.HostAPI + ":" + cfg.MinIO.PortAPI
	minioClient, err := storage.NewMinIOClient(
		minioEndpoint,
		cfg.MinIO.User,
		cfg.MinIO.Password,
		false, // useSSL = false for local development
	)
	if err != nil {
		log.Fatal("failed to create MinIO client:", err)
	}
	log.Println("✅ MinIO client initialized")

	// Ensure bucket exists
	ctx := context.Background()
	bucketName := "nexus-files"
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal("failed to check bucket:", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName)
		if err != nil {
			log.Fatal("failed to create bucket:", err)
		}
		log.Println("✅ Created MinIO bucket:", bucketName)
	} else {
		log.Println("✅ MinIO bucket exists:", bucketName)
	}

	// Step 5: AI Worker Client作成
	aiClient := client.NewAIWorkerClient("http://localhost:8001")
	log.Println("✅ AI Worker client created")

	// Step 6: Qdrant Client作成
	qdrantClient := client.NewQdrantClient("http://localhost:6333")
	log.Println("✅ Qdrant client created")

	// Step 4: File Service作成
	fileService := service.NewFileService(queries, minioClient, bucketName, qdrantClient)
	log.Println("✅ File service created")

	// Step8: OllamaClient ChatService のインスタンス作成
	ollamaBaseURL := fmt.Sprintf("http://%s:%s", cfg.Ollama.Host, cfg.Ollama.Port)
	ollamaModel := cfg.Ollama.Model
	ollamaClient := client.NewOllamaClient(ollamaBaseURL)
	log.Println("✅ Ollama client created")

	// Step 7: Document Processor作成
	documentProcessor := service.NewDocumentProcessor(queries, aiClient, qdrantClient, minioClient, ollamaClient)
	log.Println("✅ Document processor created")

	warmupCtx, warmupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer warmupCancel()

	if err := ollamaClient.WarmUp(warmupCtx, ollamaModel); err != nil {
		log.Printf("⚠️  Ollama warmup failed (non-fatal): %v", err)
		log.Println("⚠️  First chat request may be slow or fail. Consider checking Ollama service.")
		// 本番環境では fatal にする場合:
		// log.Fatalf("❌ Ollama warmup failed: %v", err)
	}

	chatService := service.NewChatService(queries, aiClient, qdrantClient, ollamaClient)
	log.Println("✅ Chat service created")

	searchService := service.NewSearchService(
		queries,
		aiClient,
		qdrantClient,
	)
	log.Println("✅ Search service created")

	analysisService := service.NewAnalysisService(queries, aiClient, qdrantClient, ollamaClient)
	log.Println("✅ Analysis service created")

	sourceService := service.NewSourceService(queries)
	log.Println("✅ Source service created")

	// --- Handler ---
	h := handler.NewHandler(database, fileService, documentProcessor, searchService, chatService, analysisService, sourceService)
	log.Println("✅ Handler created")

	// --- Router Setup ---
	r := chi.NewRouter()

	r.Use(middleware.CORS)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)

	api.HandlerWithOptions(h, api.ChiServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: r,
	})

	// --- Start Server ---
	port := cfg.Server.Port
	log.Println("🚀 Server starting on :" + port)
	log.Println("📁 File upload: POST http://localhost:8080/api/v1/workspaces/{id}/files/upload")
	log.Println("🔍 RAG search: POST http://localhost:" + port + "/api/v1/workspaces/{id}/search")
	log.Println("💬 Chat: POST http://localhost:" + port + "/api/v1/workspaces/{id}/chats/{chatId}/messages")
	log.Println("💬 Health check: GET http://localhost:8080/api/v1/health")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

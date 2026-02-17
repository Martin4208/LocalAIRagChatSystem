package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/service"
)

func main() {
	ctx := context.Background()

	// Step 1: データベース接続
	connStr := "host=localhost port=5432 user=user password=password dbname=nexus sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")

	queries := db.New(database)

	// Step 2: AI Worker Client作成
	aiClient := client.NewAIWorkerClient("http://localhost:8001")
	fmt.Println("✅ AI Worker Client created")

	// Step 3: Qdrant Client作成
	qdrantClient := client.NewQdrantClient("http://localhost:6333")
	fmt.Println("✅ Qdrant Client created")

	// Step 4: Document Processor作成
	processor := service.NewDocumentProcessor(queries, aiClient, qdrantClient)
	fmt.Println("✅ Document Processor created: %T\n", processor)

	// Step 5: テスト用のWorkspaceを作成
	workspace, err := queries.CreateWorkspace(ctx, db.CreateWorkspaceParams{
		Name:        "Test Workspace",
		Description: sql.NullString{String: "For testing document processor", Valid: true},
	})
	if err != nil {
		log.Fatalf("Failed to create workspace: %v", err)
	}
	fmt.Printf("✅ Created workspace: %s\n", workspace.ID)

	// Step 6: テスト用のテキストファイルを作成
	testContent := `This is a test document about machine learning.
Machine learning is a subset of artificial intelligence.
It enables computers to learn from data without being explicitly programmed.
Deep learning is a type of machine learning that uses neural networks.
Neural networks are inspired by the human brain.`

	testFilePath := filepath.Join(os.TempDir(), "test_document.txt")
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFilePath)
	fmt.Println("✅ Created test file")

	// Step 7: MinIOにアップロード（簡易版：手動でFileとDocumentレコード作成）
	// 注: 実際のアップロードロジックは後で実装するので、ここでは手動で作成

	// TODO: 実際にはFileServiceを使ってMinIOにアップロードする
	// 今回は簡略化のため、既存のファイルがあると仮定してテスト
	fmt.Println("⚠️  Note: このテストはMinIOに既存のファイルがあることを前提としています")
	fmt.Println("⚠️  実際のE2Eテストは後で実装します")

	// ダミーのDocumentIDでテスト
	// 実際には存在するdocument_idを指定する必要がある
	fmt.Println("\n=== このテストを完全に実行するには ===")
	fmt.Println("1. 実際にファイルをアップロードするか")
	fmt.Println("2. 既存のdocument_idを指定する必要があります")
	fmt.Println("\n今回は構造の確認のみを行いました。")
	fmt.Println("次のステップでFile Upload Handlerを実装すれば、完全なE2Eテストが可能になります。")
}

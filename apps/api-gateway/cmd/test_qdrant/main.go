package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/client"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	// Qdrant Clientを作成
	qdrantClient := client.NewQdrantClient("http://localhost:6333")

	fmt.Println("=== Step 1: Collection作成 ===")
	collectionName := "test_embeddings"
	vectorSize := 1024 // multilingual-e5-largeの次元数

	err := qdrantClient.EnsureCollection(ctx, collectionName, vectorSize)
	if err != nil {
		log.Fatalf("Failed to ensure collection: %v", err)
	}
	fmt.Printf("✅ Collection '%s' created (or already exists)\n\n", collectionName)

	// Step 2: テストデータを保存
	fmt.Println("=== Step 2: データ保存 ===")

	// ダミーのEmbeddingベクトル（実際にはAI Workerから取得）
	dummyVector1 := make([]float64, vectorSize)
	dummyVector2 := make([]float64, vectorSize)
	for i := 0; i < vectorSize; i++ {
		dummyVector1[i] = float64(i) / 1000.0
		dummyVector2[i] = float64(i+1) / 1000.0
	}

	points := []client.Point{
		{
			ID:     uuid.New().String(),
			Vector: dummyVector1,
			Payload: map[string]interface{}{
				"document_id": "doc-123",
				"chunk_index": 0,
				"text":        "This is the first test chunk about machine learning",
			},
		},
		{
			ID:     uuid.New().String(),
			Vector: dummyVector2,
			Payload: map[string]interface{}{
				"document_id": "doc-123",
				"chunk_index": 1,
				"text":        "This is the second test chunk about artificial intelligence",
			},
		},
	}

	err = qdrantClient.UpsertPoints(ctx, collectionName, points)
	if err != nil {
		log.Fatalf("Failed to upsert points: %v", err)
	}
	fmt.Printf("✅ Saved %d points to Qdrant\n\n", len(points))

	// Step 3: 検索テスト
	fmt.Println("=== Step 3: 検索テスト ===")

	// dummyVector1に近いベクトルを検索
	searchVector := dummyVector1

	time.Sleep(1 * time.Second) // Qdrantがインデックスを更新するまで少し待つ

	results, err := qdrantClient.Search(ctx, collectionName, searchVector, 5)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("Found %d results:\n", len(results.Result))
	for i, result := range results.Result {
		text := result.Payload["text"]
		fmt.Printf("%d. Score: %.4f, Text: %v\n", i+1, result.Score, text)
	}

	fmt.Println("\n✅ All tests passed!")
}

# Nexus Local - 開発TODO

## 🎯 現在の状態
- [x] Docker環境構築（MinIO, Qdrant, PostgreSQL）
- [x] Next.jsフロントエンド起動確認
- [x] 要件定義・ER図作成
- [ ] バックエンド開発（← **今ここ**）

---

## 📋 Phase 1: データベース基盤構築

### 1.1 PostgreSQL マイグレーション
- [ ] `migrations/` ディレクトリ作成
- [ ] `001_create_documents.sql` 作成
- [ ] `002_create_chat.sql` 作成
- [ ] `003_create_analysis.sql` 作成
- [ ] マイグレーション実行コマンドを `Taskfile.yml` に追加

### 1.2 テストデータ投入
- [ ] サンプルファイル（PDF, 画像）を準備
- [ ] `seeds/` ディレクトリ作成
- [ ] 初期データ投入スクリプト

---

## 📋 Phase 2: バックエンドAPI（最小構成）

### 2.1 Go API Gateway セットアップ
- [ ] `apps/api-gateway/` に main.go 作成
- [ ] ヘルスチェックエンドポイント（GET /health）
- [ ] CORS設定（Next.jsから叩けるように）

### 2.2 ファイルアップロード機能
- [ ] `POST /api/upload` エンドポイント
- [ ] MinIOへのファイル保存処理
- [ ] PostgreSQLへのメタデータ記録
- [ ] エラーハンドリング

### 2.3 ファイル一覧取得
- [ ] `GET /api/documents` エンドポイント
- [ ] PostgreSQLからデータ取得
- [ ] JSON整形して返す

---

## 📋 Phase 3: フロントエンド統合

### 3.1 アップロード画面
- [ ] ファイル選択UI（shadcn/ui使用）
- [ ] アップロードボタン
- [ ] プログレス表示
- [ ] 成功・失敗通知

### 3.2 ファイル一覧画面
- [ ] テーブル表示
- [ ] サムネイル表示（画像の場合）
- [ ] 削除ボタン

---

## 📋 Phase 4: AI機能（MVP）

### 4.1 テキスト抽出
- [ ] Python Worker セットアップ
- [ ] PDFからテキスト抽出
- [ ] 画像OCR（Tesseract）
- [ ] 抽出結果をDBに保存

### 4.2 Embedding生成
- [ ] Sentence Transformersセットアップ
- [ ] テキストをベクトル化
- [ ] Qdrantに保存

### 4.3 基本的なRAG検索
- [ ] `POST /api/search` エンドポイント
- [ ] Qdrantでベクトル検索
- [ ] 関連ドキュメント返却

---

## 📋 Phase 5: チャット機能

- [ ] `POST /api/chat` エンドポイント
- [ ] Ollamaとの連携
- [ ] RAG結果をコンテキストに含める
- [ ] ストリーミングレスポンス（将来）

---

## 🚧 保留・将来対応
- [ ] Neo4j グラフDB統合
- [ ] 分析機能
- [ ] 動画処理（Whisper）
- [ ] gRPC化（必要になったら）
- [ ] Kubernetes移行

---

## 💡 メモ・気づき

### 2024-12-17
- gRPCは個人開発では過剰 → REST APIで進める
- まずはファイルアップロード1機能を完全に動かすことに集中

---

## 🔗 関連ファイル
- 要件定義: `/mnt/project/要件定義`
- ER図: `/mnt/project/Nexus_ER図.pdf`
- Taskfile: `/mnt/project/Taskfile.yml`

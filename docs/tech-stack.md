## 技術スタック

### フロントエンド
- Next.js(App Router + Static Export)
- TypeScript
- Tailwind CSS
- shadcn/ui

### バックエンド
- API Gateway : Go (REST)
- AI Worker : Python (Embedding, Whisper, CLIP)
- Communication : Protocol Buffers (gRPC)

### インフラ
- Object Storage : MinIO (S3-compatible)
- Metadata DB : PostgreSQL
- Vector DB : Qdrant
- Graph DB : Neo4j

- Containerization : Docker 
- Container Orchestration : Docker Compose (future: Kubernetes)
- Workflow Orchestration :Apache Airflow or Dagster

- Desktop Runtime : Tauri

- LLM Runtime : Ollama (Llama 3.2, Qwen)
- Sentence Transformers (Embedding)
- Whisper (音声文字起こし)

### 監視・運用
- Observability: OpenTelemetry + Grafana

### 開発ツール
- IDE : VSCode
- Task Runner : Taskfile (task)
<!-- IaC : Terraform or Ansible -->
- Version Control : Git 
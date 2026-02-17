


Local Setup
bash# インフラ起動
docker compose up -d

# API Gateway
cd apps/api-gateway
go run ./cmd/server

# AI Worker
cd apps/ai-worker
pip install -r requirements.txt
uvicorn main:app --port 8001

# Frontend
cd apps/web
npm install && npm run dev

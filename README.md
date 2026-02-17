## Unfinished project
### This is connected to ollama, using phi3:mini. You can upload pdf or txt files and the RAG system will let phi3:mini to respond from the documents.


Local Setup
docker compose up -d

API Gateway
cd apps/api-gateway
go run ./cmd/server

AI Worker
cd apps/ai-worker
pip install -r requirements.txt
uvicorn main:app --port 8001

Frontend
cd apps/web
npm install && npm run dev

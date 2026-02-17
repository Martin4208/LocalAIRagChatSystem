# アプリケーションのエントリーポイント (=>API)
import logging
from fastapi import FastAPI
from routers import embedding, embed_query, generate
from api.deps import get_model
from utils.logging import setup_logging

logger = logging.getLogger(__name__)

app = FastAPI(
    title="Nexus AI Worker",
    description="Embedding and AI processing service",
    version="1.0.0"
)

app.include_router(
    embedding.router,
    prefix="/api/v1",
    tags=["embedding"]
)

app.include_router(
    embed_query.router,
    prefix="/api/v1",
    tags=["query"]
)

app.include_router(
    generate.router,
    prefix="/api/v1",
    tags=["generate"]
)

@app.get("/health")
async def health():
    return {"status": "ok"}

@app.on_event("startup")
async def startup_event():
    """
    アプリ起動時に事前ロード
    """
    setup_logging()
    
    logger.info("Loading embedding model...")
    model = get_model()
    
    model_info = model.get_model_info()
    logger.info(f"Model loaded: {model_info['model_name']}")
    logger.info(f"Dimension: {model_info['embedding_dim']}")
    logger.info(f"AI Worker is ready!")

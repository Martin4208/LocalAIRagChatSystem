from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from models.embedding import get_embedding_model

router = APIRouter()

class EmbedQueryRequest(BaseModel):
    text: str
    
class EmbedQueryResponse(BaseModel):
    embedding: list[float]
    dim: int

@router.post("/embed/query", response_model=EmbedQueryResponse)
async def embed_query(request: EmbedQueryRequest):
    try:
        model = get_embedding_model()
        
        embedding = model.encode_query(request.text)
        
        return EmbedQueryResponse(
            embedding=embedding.tolist(),
            dim=len(embedding)
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
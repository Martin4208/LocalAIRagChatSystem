from fastapi import APIRouter, Depends, HTTPException
from schemas.embedding import EmbedRequest, EmbedResponse, ErrorResponse
from models.embedding import EmbeddingModel
from api.deps import get_model
import time
from exceptions import ModelLoadError

router = APIRouter()

@router.post("/embed", response_model=EmbedResponse)
async def embed_texts(
    request: EmbedRequest,
    model: EmbeddingModel = Depends(get_model)
):
    """
    テキストをEmbedding化するエンドポイント

    Args:
        request (EmbedRequest): _description_
        model (EmbeddingModel, optional): _description_. Defaults to Depends(get_model).
    """
    try:
        start = time.perf_counter()
        
        embeddings_np = model.encode_documents(request.texts) # shape: (N, 1024) numpy array
        model_info = model.get_model_info()
        embeddings_list = embeddings_np.tolist()
        
        elapsed = (time.perf_counter() - start) * 1000
        
        return EmbedResponse(
            embeddings=embeddings_list,
            model=model_info["model_name"],
            dimension=model_info["embedding_dim"],
            count=len(request.texts),
            elapsed_ms=elapsed
        )
        
    except ModelLoadError as e:
        # Model特有のエラー → 503
        raise HTTPException(
            status_code=503,
            detail="AI model is not available"
        )
        
    except ValueError as e:
        # データの問題 → 400
        raise HTTPException(
            status_code=400,
            detail=f"Invalid input: {str(e)}"
        )
        
    except Exception as e:
        # その他の予期しないエラー → 500
        raise HTTPException(
            status_code=500,
            detail="Internal processing error"
        )

        
# Dependency Injection = (必要なものを外から注入する設計パターン)
from functools import lru_cache
from models.embedding import get_embedding_model, EmbeddingModel

@lru_cache()
def get_model() -> EmbeddingModel:
    """
    アプリケーション起動後、最初の呼び出し時に1回だけ実行
    ２回目以降はキャッシュされたインスタンスを返す

    Returns:
        EmbeddingModel:
    """
    return get_embedding_model()
    

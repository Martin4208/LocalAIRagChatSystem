import logging
from sentence_transformers import SentenceTransformer
from typing import Optional
from functools import lru_cache
import numpy as np
from config import get_settings
from exceptions import ModelLoadError

logger = logging.getLogger(__name__)

class EmbeddingModel():
    def __init__(self):
        self.settings = get_settings()
        self.model: Optional[SentenceTransformer] = None
        self.model_name: str = self.settings.embedding_model
        self.embedding_dim: int = self.settings.embedding_dim
        self.is_loaded: bool = False
        
        self._load_model()
    
    def encode_documents(self, texts: list[str]) -> np.ndarray:
        """
        Encode documents into embeddings.

        Args:
            texts (list[str]): List of document texts

        Returns:
            np.ndarray: Shape (len(texts), embedding_dim)
        """
        if not self.is_loaded or self.model is None:
            raise RuntimeError("Model is not loaded")
        texts = [f"passage: {t}" for t in texts]
        return self.model.encode(texts, batch_size=self.settings.max_batch_size, convert_to_numpy=True, normalize_embeddings=True)
        
    def encode_query(self, text: str) -> np.ndarray:
        if not self.is_loaded or self.model is None:
            raise RuntimeError("Model is not loaded")
        text = f"query: {text}"
        return self.model.encode(text, convert_to_numpy=True, normalize_embeddings=True)
        
    
    # Utility Methods
    def get_model_info(self) -> dict:
        return {
            "model_name": self.model_name,
            "embedding_dim": self.embedding_dim,
            "is_loaded": self.is_loaded
        }
    
    def is_model_loaded(self) -> bool:
        return self.is_loaded
    
    def _load_model(self):
        try:
            self.model = SentenceTransformer(
                self.model_name,
                cache_folder=self.settings.model_cache_dir
            )
            self.is_loaded = True
            logger.info(f"Model loaded: {self.model_name}")
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            raise ModelLoadError(f"Cannot load {self.model_name}") from e


@lru_cache()
def get_embedding_model() -> EmbeddingModel:
    return EmbeddingModel()
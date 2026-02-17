from pydantic import BaseModel, Field
from enum import Enum
from typing import Optional

class EmbedRequest(BaseModel):
    texts: list[str]

class EmbedResponse(BaseModel):
    embeddings: list[list[float]]
    count: int
    model: Optional[str]
    dimension: Optional[int]
    elapsed_ms: Optional[float]


class ErrorCode(str, Enum):
    MODEL_NOT_LOADED = "MODEL_NOT_LOADED"
    INVALID_INPUT = "INVALID_INPUT"
    PROCESSING_ERROR = "PROCESSING_ERROR"


class ErrorResponse(BaseModel):
    code: ErrorCode
    message: str
    details: Optional[dict] = None
# apps/ai-worker/schemas/__init__.py

"""
Schemas module

API request/response schemas for FastAPI
"""

from .requests import (
    EmbedDocumentsRequest,
    EmbedDocumentsResponse,
    EmbedQueryRequest,
    EmbedQueryResponse,
    ErrorResponse,
    HealthResponse,
)

__all__ = [
    "EmbedDocumentsRequest",
    "EmbedDocumentsResponse",
    "EmbedQueryRequest",
    "EmbedQueryResponse",
    "HealthResponse",
    "ErrorResponse",
]
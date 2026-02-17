# apps/ai-worker/schemas/requests.py

"""
API Request/Response Schemas

FastAPI で使用する Pydantic スキーマ定義。
自動的に以下を提供:
- リクエストバリデーション
- レスポンスのシリアライゼーション
- OpenAPI (Swagger) ドキュメント生成
- 型ヒント（VSCode補完）
"""

from typing import Literal, Optional

from pydantic import BaseModel, Field, field_validator


# ========================================
# POST /embed/documents
# ========================================

class EmbedDocumentsRequest(BaseModel):
    """
    複数ドキュメントのEmbeddingリクエスト
    
    ドキュメント保存時に使用。最大100個まで同時処理可能。
    
    Attributes:
        texts: Embedding対象のテキストリスト（1-100個）
    
    Examples:
        >>> request = EmbedDocumentsRequest(
        ...     texts=["First document", "Second document"]
        ... )
        >>> len(request.texts)
        2
    
    Raises:
        ValidationError: textsが空、または空文字列を含む場合
    """
    
    texts: list[str] = Field(
        ...,
        min_length=1,
        max_length=100,
        description="List of texts to embed (1-100 items)",
        examples=[
            ["This is the first document", "This is the second document"]
        ]
    )
    
    @field_validator('texts')
    @classmethod
    def validate_texts(cls, v: list[str]) -> list[str]:
        """
        テキストリストのバリデーション
        
        チェック内容:
        1. 空文字列または空白のみの要素がないか
        2. 各テキストが最大文字数を超えていないか
        
        Args:
            v: バリデーション対象のテキストリスト
            
        Returns:
            バリデーション済みのテキストリスト
            
        Raises:
            ValueError: バリデーションエラー時
        """
        MAX_TEXT_LENGTH = 10000
        
        # 空文字列・空白のみチェック
        for i, text in enumerate(v):
            if not text or not text.strip():
                raise ValueError(
                    f"texts[{i}] is empty or contains only whitespace"
                )
        
        # 文字数制限チェック
        for i, text in enumerate(v):
            if len(text) > MAX_TEXT_LENGTH:
                raise ValueError(
                    f"texts[{i}] exceeds maximum length of {MAX_TEXT_LENGTH} "
                    f"characters (got {len(text)})"
                )
        
        return v
    
    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "texts": [
                        "The quick brown fox jumps over the lazy dog",
                        "Machine learning is a subset of artificial intelligence"
                    ]
                }
            ]
        }
    }


class EmbedDocumentsResponse(BaseModel):
    """
    複数Embeddingのレスポンス
    
    Attributes:
        embeddings: Embeddingベクトルのリスト（外側: テキスト数、内側: 次元数）
        count: 生成されたEmbeddingの数
        model: 使用したモデル名（オプション）
        dim: ベクトルの次元数（オプション）
        processing_time: 処理時間（秒、オプション）
    
    Examples:
        >>> response = EmbedDocumentsResponse(
        ...     embeddings=[[0.1, 0.2, 0.3], [0.4, 0.5, 0.6]],
        ...     count=2,
        ...     dim=3
        ... )
        >>> response.count
        2
    """
    
    embeddings: list[list[float]] = Field(
        ...,
        description="List of embedding vectors (outer: texts, inner: dimensions)"
    )
    
    count: int = Field(
        ...,
        ge=0,
        description="Number of embeddings generated"
    )
    
    model: Optional[str] = Field(
        default=None,
        description="Model used for embedding generation"
    )
    
    dim: Optional[int] = Field(
        default=None,
        ge=1,
        description="Dimension of each embedding vector"
    )
    
    processing_time: Optional[float] = Field(
        default=None,
        ge=0,
        description="Processing time in seconds"
    )
    
    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "embeddings": [
                        [0.123, -0.456, 0.789],
                        [-0.234, 0.567, -0.890]
                    ],
                    "count": 2,
                    "model": "intfloat/multilingual-e5-large",
                    "dim": 3,
                    "processing_time": 0.152
                }
            ]
        }
    }


# ========================================
# POST /embed/query
# ========================================

class EmbedQueryRequest(BaseModel):
    """
    単一クエリのEmbeddingリクエスト
    
    検索クエリの処理に使用。1つのテキストのみを受け付ける。
    
    Attributes:
        text: Embedding対象のクエリテキスト
    
    Examples:
        >>> request = EmbedQueryRequest(text="What is machine learning?")
        >>> request.text
        'What is machine learning?'
    
    Raises:
        ValidationError: textが空、または長すぎる場合
    """
    
    text: str = Field(
        ...,
        min_length=1,
        max_length=10000,
        description="Query text to embed (single text only)"
    )
    
    @field_validator('text')
    @classmethod
    def validate_text(cls, v: str) -> str:
        """
        クエリテキストのバリデーション
        
        空白のみの文字列を弾く
        
        Args:
            v: バリデーション対象のテキスト
            
        Returns:
            バリデーション済みのテキスト
            
        Raises:
            ValueError: テキストが空白のみの場合
        """
        if not v.strip():
            raise ValueError("text cannot be empty or whitespace only")
        return v
    
    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "text": "What are the risks of budget cuts?"
                }
            ]
        }
    }


class EmbedQueryResponse(BaseModel):
    """
    単一Embeddingのレスポンス
    
    Attributes:
        embedding: Embeddingベクトル
        dim: ベクトルの次元数
        model: 使用したモデル名（オプション）
        processing_time: 処理時間（秒、オプション）
    
    Examples:
        >>> response = EmbedQueryResponse(
        ...     embedding=[0.1, 0.2, 0.3],
        ...     dim=3
        ... )
        >>> len(response.embedding)
        3
    """
    
    embedding: list[float] = Field(
        ...,
        description="Embedding vector for the query"
    )
    
    dim: int = Field(
        ...,
        ge=1,
        description="Dimension of the embedding vector"
    )
    
    model: Optional[str] = Field(
        default=None,
        description="Model used for embedding generation"
    )
    
    processing_time: Optional[float] = Field(
        default=None,
        ge=0,
        description="Processing time in seconds"
    )
    
    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "embedding": [0.123, -0.456, 0.789],
                    "dim": 3,
                    "model": "intfloat/multilingual-e5-large",
                    "processing_time": 0.045
                }
            ]
        }
    }


# ========================================
# GET /health
# ========================================

class HealthResponse(BaseModel):
    """
    ヘルスチェックレスポンス
    
    サービスの状態とモデルのロード状況を返す
    
    Attributes:
        status: サービスの状態（healthy/unhealthy）
        model_loaded: モデルがメモリにロード済みか
        model_name: 使用中のモデル名
        uptime_seconds: サービス起動からの経過時間（秒）
        memory_usage_mb: メモリ使用量（MB、オプション）
    
    Examples:
        >>> response = HealthResponse(
        ...     status="healthy",
        ...     model_loaded=True,
        ...     model_name="intfloat/multilingual-e5-large",
        ...     uptime_seconds=3600.5
        ... )
        >>> response.status
        'healthy'
    """
    
    status: Literal["healthy", "unhealthy"] = Field(
        ...,
        description="Service health status"
    )
    
    model_loaded: bool = Field(
        ...,
        description="Whether the embedding model is loaded in memory"
    )
    
    model_name: str = Field(
        ...,
        description="Name of the loaded model"
    )
    
    uptime_seconds: float = Field(
        ...,
        ge=0,
        description="Service uptime in seconds"
    )
    
    memory_usage_mb: Optional[float] = Field(
        default=None,
        ge=0,
        description="Current memory usage in megabytes"
    )
    
    model_config = {
        "protected_namespaces": (),
        "json_schema_extra": {
            "examples": [
                {
                    "status": "healthy",
                    "model_loaded": True,
                    "model_name": "intfloat/multilingual-e5-large",
                    "uptime_seconds": 3600.5,
                    "memory_usage_mb": 2048.3
                }
            ]
        }
    }


# ========================================
# Error Response (Common)
# ========================================

class ErrorResponse(BaseModel):
    """
    エラーレスポンス（共通）
    
    全てのエラーで統一されたフォーマットを使用
    
    Attributes:
        code: エラーコード（例: VALIDATION_ERROR, MODEL_ERROR）
        message: エラーメッセージ（人間が読める形式）
        details: 詳細情報（オプション）
    
    Examples:
        >>> error = ErrorResponse(
        ...     code="VALIDATION_ERROR",
        ...     message="Invalid input format",
        ...     details={"field": "texts", "issue": "empty_list"}
        ... )
        >>> error.code
        'VALIDATION_ERROR'
    """
    
    code: str = Field(
        ...,
        description="Error code for programmatic handling"
    )
    
    message: str = Field(
        ...,
        description="Human-readable error message"
    )
    
    details: Optional[dict] = Field(
        default=None,
        description="Additional error details (optional)"
    )
    
    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "code": "VALIDATION_ERROR",
                    "message": "Invalid input: texts cannot be empty",
                    "details": {
                        "field": "texts",
                        "error": "empty_list",
                        "min_length": 1
                    }
                },
                {
                    "code": "MODEL_ERROR",
                    "message": "Failed to generate embeddings",
                    "details": {
                        "reason": "model_not_loaded"
                    }
                }
            ]
        }
    }


# ========================================
# __init__.py で export するための定義
# ========================================

__all__ = [
    "EmbedDocumentsRequest",
    "EmbedDocumentsResponse",
    "EmbedQueryRequest",
    "EmbedQueryResponse",
    "HealthResponse",
    "ErrorResponse",
]
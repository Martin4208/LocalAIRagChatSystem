# apps/ai-worker/config.py

"""
AI Worker Configuration

環境変数または .env ファイルから設定を読み込み、
型安全性とバリデーションを提供します。

使用例:
    from config import get_settings
    
    settings = get_settings()
    print(settings.embedding_model)
"""

from functools import lru_cache
from typing import Literal, Optional

from pydantic import Field, field_validator
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """
    アプリケーション設定
    
    全ての設定項目は環境変数または .env ファイルから読み込み可能。
    環境変数名は大文字小文字を区別しない（case_sensitive=False）。
    
    例:
        EMBEDDING_MODEL=my-model
        embedding_model=my-model
        どちらも同じ扱い
    """
    
    # ========================================
    # Model Settings
    # ========================================
    
    embedding_model: str = Field(
        default="intfloat/multilingual-e5-large",
        description="Hugging Face model name or local path"
    )
    """
    使用するEmbeddingモデル
    - Hugging Faceのモデル名（例: intfloat/multilingual-e5-large）
    - ローカルパス（例: ./models/my-model）
    """
    
    embedding_dim: int = Field(
        default=1024,
        ge=1,
        description="Embedding vector dimension"
    )
    """
    Embeddingベクトルの次元数
    - multilingual-e5-large: 1024
    - multilingual-e5-base: 768
    モデルと一致させる必要がある
    """
    
    max_batch_size: int = Field(
        default=32,
        ge=1,
        le=128,
        description="Maximum batch size for encoding"
    )
    """
    一度に処理する最大テキスト数
    - 大きい: 処理速度up、メモリ使用量up
    - 小さい: メモリ安全、レイテンシdown
    CPU環境推奨値: 16-32
    """
    
    # ========================================
    # Server Settings
    # ========================================
    
    server_host: str = Field(
        default="0.0.0.0",
        description="Server bind address"
    )
    """
    サーバーのバインドアドレス
    - 0.0.0.0: 全てのインターフェースで待ち受け（Docker推奨）
    - 127.0.0.1: ローカルホストのみ
    """
    
    server_port: int = Field(
        default=8001,
        ge=1024,
        le=65535,
        description="Server port number"
    )
    """
    サーバーのポート番号
    - 1024-65535: 非特権ポート範囲
    - 8001: デフォルト（Go APIは8080）
    """
    
    reload: bool = Field(
        default=False,
        description="Enable auto-reload on code changes"
    )
    """
    コード変更時の自動リロード
    - True: 開発時に便利（--reload）
    - False: 本番環境では必ずFalse
    """
    
    # ========================================
    # Performance Settings
    # ========================================
    
    num_workers: int = Field(
        default=1,
        ge=1,
        le=4,
        description="Number of worker processes"
    )
    """
    ワーカープロセス数
    - 1: モデルを1つだけメモリに載せる（推奨）
    - 2+: モデルが複数コピーされる（メモリ使用量 × N）
    
    注意: CPU環境では1推奨。GPU環境でも1で十分。
    """
    
    request_timeout: int = Field(
        default=30,
        ge=1,
        description="Request timeout in seconds"
    )
    """
    リクエストタイムアウト（秒）
    - 通常のEmbedding処理: 1-5秒
    - 大量バッチ処理: 10-30秒
    """
    
    model_load_timeout: int = Field(
        default=300,
        ge=60,
        description="Model loading timeout in seconds"
    )
    """
    モデルロードのタイムアウト（秒）
    - 初回ダウンロード: 2-5分（モデルサイズによる）
    - 2回目以降（キャッシュあり）: 10-30秒
    """
    
    # ========================================
    # Logging Settings
    # ========================================
    
    log_level: Literal["DEBUG", "INFO", "WARNING", "ERROR"] = Field(
        default="INFO",
        description="Logging level"
    )
    """
    ログレベル
    - DEBUG: 詳細なデバッグ情報
    - INFO: 通常の動作ログ（推奨）
    - WARNING: 警告のみ
    - ERROR: エラーのみ
    """
    
    log_format: Literal["text", "json"] = Field(
        default="json",
        description="Log output format"
    )
    """
    ログ出力形式
    - text: 人間が読みやすい形式
    - json: 構造化ログ（Grafana/Lokiで解析可能）
    """
    
    # ========================================
    # Cache Settings
    # ========================================
    
    model_cache_dir: str = Field(
        default="./models_cache",
        description="Directory to cache downloaded models"
    )
    """
    モデルファイルのキャッシュディレクトリ
    - 初回起動時にモデルをダウンロード
    - 2回目以降はこのディレクトリから読み込み
    - Docker volumeでマウント推奨
    """
    
    # ========================================
    # Optional Settings
    # ========================================
    
    max_memory_gb: Optional[float] = Field(
        default=None,
        ge=0.1,
        description="Memory usage warning threshold in GB"
    )
    """
    メモリ使用量の警告閾値（GB）
    - None: チェックしない
    - 8.0: 8GB超えたら警告ログ
    
    注意: モデルだけで2GB程度使用
    """
    
    enable_health_check: bool = Field(
        default=True,
        description="Enable /health endpoint"
    )
    """
    Health checkエンドポイントの有効化
    - True: GET /health が使える
    - False: 無効化（セキュリティ理由）
    """
    
    # ========================================
    # Pydantic Configuration
    # ========================================
    
    model_config = {
        "env_file": ".env",
        "env_file_encoding": "utf-8",
        "case_sensitive": False,
        "extra": "ignore",  # 未知の環境変数を無視
    }
    """
    Pydantic設定
    - env_file: .envファイルから読み込み（なくてもOK）
    - case_sensitive: 環境変数の大文字小文字を区別しない
    - extra: 未定義の環境変数があってもエラーにしない
    """
    
    # ========================================
    # Custom Validators
    # ========================================
    
    @field_validator("model_cache_dir")
    @classmethod
    def validate_cache_dir(cls, v: str) -> str:
        """
        キャッシュディレクトリのバリデーション
        
        相対パスを絶対パスに変換し、必要に応じてディレクトリを作成
        """
        from pathlib import Path
        
        path = Path(v)
        
        # 相対パスなら絶対パスに変換
        if not path.is_absolute():
            path = path.resolve()
        
        # ディレクトリが存在しない場合は作成
        path.mkdir(parents=True, exist_ok=True)
        
        return str(path)
    
    @field_validator("num_workers")
    @classmethod
    def validate_num_workers(cls, v: int) -> int:
        """
        ワーカー数のバリデーション
        
        CPU数を超える場合は警告を出す
        """
        import os
        
        cpu_count = os.cpu_count() or 1
        
        if v > cpu_count:
            import warnings
            warnings.warn(
                f"num_workers ({v}) exceeds CPU count ({cpu_count}). "
                f"This may cause performance degradation.",
                UserWarning
            )
        
        return v


@lru_cache()
def get_settings() -> Settings:
    """
    設定インスタンスを取得（シングルトン）
    
    @lru_cache() により、アプリケーション起動中は
    常に同じSettingsインスタンスを返します。
    
    これにより:
    - .envファイルの再読み込みを防ぐ
    - メモリ効率が良い
    - 設定の一貫性を保証
    
    Returns:
        Settings: 設定インスタンス
    
    Example:
        >>> from config import get_settings
        >>> settings = get_settings()
        >>> print(settings.embedding_model)
        intfloat/multilingual-e5-large
    """
    return Settings()


# ========================================
# Convenience Functions
# ========================================

def get_model_name() -> str:
    """モデル名を取得するショートカット"""
    return get_settings().embedding_model


def get_batch_size() -> int:
    """バッチサイズを取得するショートカット"""
    return get_settings().max_batch_size


def is_development() -> bool:
    """開発環境かどうかを判定"""
    return get_settings().reload


# ========================================
# For Testing
# ========================================

def reset_settings_cache() -> None:
    """
    設定キャッシュをクリア（テスト用）
    
    通常のアプリケーションでは使用しない。
    ユニットテストで環境変数を変更した後に呼ぶ。
    """
    get_settings.cache_clear()


if __name__ == "__main__":
    # 設定の確認用スクリプト
    settings = get_settings()
    
    print("=" * 50)
    print("AI Worker Configuration")
    print("=" * 50)
    print(f"Model: {settings.embedding_model}")
    print(f"Dimension: {settings.embedding_dim}")
    print(f"Batch Size: {settings.max_batch_size}")
    print(f"Server: {settings.server_host}:{settings.server_port}")
    print(f"Workers: {settings.num_workers}")
    print(f"Log Level: {settings.log_level}")
    print(f"Cache Dir: {settings.model_cache_dir}")
    print("=" * 50)
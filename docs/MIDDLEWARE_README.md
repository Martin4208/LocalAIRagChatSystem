# Middleware Package

このパッケージには、Enterprise グレードの HTTP ミドルウェアが含まれています。

## 📦 含まれるミドルウェア

### 1. Recovery Middleware
パニックをキャッチし、サーバーがクラッシュするのを防ぎます。

**機能:**
- パニックを捕捉
- スタックトレースをログに記録
- 500 Internal Server Error を返却
- JSON形式のエラーレスポンス

### 2. Logger Middleware
すべての HTTP リクエストをログに記録します。

**ログ内容:**
- HTTPメソッド (GET, POST, etc.)
- リクエストパス
- クライアントIPアドレス
- ステータスコード
- 処理時間

**出力例:**
```
[GET] /workspaces 127.0.0.1:54321 - 200 - 3.5ms
[POST] /workspaces 127.0.0.1:54322 - 201 - 12.3ms
```

### 3. CORS Middleware
Cross-Origin Resource Sharing (CORS) を処理します。

**設定内容:**
- `Access-Control-Allow-Origin: *` (開発環境)
- `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With`
- OPTIONSプリフライトリクエストの処理

**注意:** 本番環境では `Allow-Origin` を特定のドメインに制限してください。

### 4. Chain Utility
複数のミドルウェアを簡単に適用するためのユーティリティ。

## 🚀 使い方

### 基本的な使用例

```go
import (
    "github.com/go-chi/chi/v5"
    "nexus/api-gateway/internal/middleware"
)

func main() {
    router := chi.NewRouter()
    
    // ミドルウェアを適用
    stack := middleware.Chain(
        middleware.Recovery,  // 最初に適用（最外層）
        middleware.Logger,
        middleware.CORS,
    )
    
    router.Use(stack)
    
    // ルートを定義
    router.Get("/health", HealthHandler)
}
```

### ミドルウェアの順序

**重要:** ミドルウェアの順序は重要です！

```
Request 
  → Recovery   (パニックをキャッチ)
    → Logger   (リクエストをログ)
      → CORS   (ヘッダーを設定)
        → Handler (実際の処理)
```

**推奨順序:**
1. **Recovery** - 最初に適用（すべてのパニックをキャッチ）
2. **Logger** - リクエストを記録
3. **CORS** - CORSヘッダーを設定
4. **Handler** - 実際のビジネスロジック

## 🧪 テスト方法

### 1. Logger のテスト
```bash
curl http://localhost:8080/health
# ログに [GET] /health ... が表示される
```

### 2. CORS のテスト
```bash
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/workspaces
# CORS ヘッダーが返却される
```

### 3. Recovery のテスト
```go
// テストハンドラーを追加
router.Get("/test/panic", func(w http.ResponseWriter, r *http.Request) {
    panic("test panic!")
})
```

```bash
curl http://localhost:8080/test/panic
# 500 エラーが返却され、サーバーは落ちない
# ログにパニックとスタックトレースが記録される
```

## 🏗️ カスタマイズ

### CORS を本番環境用に設定

```go
func CORSProduction(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            for _, allowed := range allowedOrigins {
                if origin == allowed {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    break
                }
            }
            
            // ... rest of CORS logic
        })
    }
}
```

### 構造化ログへの移行

```go
import "go.uber.org/zap"

func LoggerWithZap(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            lrw := newLoggingResponseWriter(w)
            
            next.ServeHTTP(lrw, r)
            
            logger.Info("request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.Int("status", lrw.statusCode),
                zap.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

## 📚 参考資料

- [Go HTTP Middleware Pattern](https://www.alexedwards.net/blog/making-and-using-middleware)
- [CORS Explained](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Effective Go - Defer, Panic, Recover](https://go.dev/doc/effective_go#recover)
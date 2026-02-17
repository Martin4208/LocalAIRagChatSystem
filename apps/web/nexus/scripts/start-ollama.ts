import { spawn } from 'child_process';
import fetch from 'node-fetch';

/**
 * Ollama をバックグラウンドで起動し、warmup メッセージを送る
 */
async function startOllama() {
  console.log('[Ollama] Starting server...');

  // Ollama サーバーを起動（変数は _ で警告回避）
  const _ollama = spawn('ollama', ['serve'], { stdio: 'inherit' });

  // サーバーが立ち上がるのを少し待ってから warmup
  const waitTime = 5000; // 5秒
  console.log(`[Ollama] Waiting ${waitTime / 1000}s for server to start...`);
  setTimeout(async () => {
    try {
      console.log('[Ollama] Sending warmup message...');
      const res = await fetch('http://127.0.0.1:11434/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: 'phi3:mini',
          prompt: 'warmup',
          max_tokens: 10
        }),
      });

      const data = await res.json();
      console.log('[Ollama] Warmup done:', data);
    } catch (err) {
      console.error('[Ollama] Warmup failed:', (err as Error).message);
    }
  }, waitTime);
}

// スクリプト単独で実行可能にする
if (require.main === module) {
  startOllama();
}

export default startOllama;

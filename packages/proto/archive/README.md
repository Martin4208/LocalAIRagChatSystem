# gRPC/Protocol Buffers アーカイブ

このディレクトリには、初期に作成したgRPC定義が保存されています。

## 経緯
- 当初gRPCで進めようとしたが、RESTの方が開発効率が良いと判断
- 将来、内部サービス間通信（API Gateway ↔ AI Worker）でgRPCが必要になったら復活

## 保存内容
- `nexus/`: Protocol Buffers定義
- `buf.yaml`, `buf.gen.yaml`: Buf設定
- `generated/`: 生成済みのGo/Python/TSコード

## 復活方法
必要になったら、これらを `packages/proto/` に戻してbuf generateを実行
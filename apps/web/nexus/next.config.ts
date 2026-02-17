import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      // MinIO (本番・開発共通)
      {
        protocol: "http",
        hostname: "localhost",
        port: "9000",
        pathname: "/**",
      },
      // 開発環境のみ：ダミー画像サービス
      ...(process.env.NODE_ENV === "development"
        ? [
            {
              protocol: "https",
              hostname: "picsum.photos",
              pathname: "/**",
            },
          ]
        : []),
    ],
  },


  /**
   * モノレポ環境で Next.js が依存解析を誤らないように
   * ワークスペースルートを明示的に指定
   *
   * apps/web/nexus → 3階層上がリポジトリルート
   */
  outputFileTracingRoot: path.join(__dirname, "../../.."),

  webpack: (config) => {
    config.resolve.alias.canvas = false;
    config.resolve.alias.encoding = false;
    return config;
  },
};

export default nextConfig;

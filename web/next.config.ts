import type { NextConfig } from "next";

const apiOrigin = process.env.CIAM_API_ORIGIN ?? "http://localhost:13702";

const nextConfig: NextConfig = {
  async rewrites() {
    return [{ source: "/api/:path*", destination: `${apiOrigin}/api/:path*` }];
  },
};

export default nextConfig;

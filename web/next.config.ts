import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    // Load images from randomuser.me for Lego driver profile pictures
    domains: ['randomuser.me'],
  },
  reactStrictMode: false,
  output: 'standalone',
};

export default nextConfig;

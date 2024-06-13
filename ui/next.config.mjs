/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    ROOT_API_URL: process.env.ROOT_API_URL,
    NEXT_PUBLIC_ROOT_API_URL: process.env.NEXT_PUBLIC_ROOT_API_URL,
  },
};

export default nextConfig;

/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    NEXT_PUBLIC_ROOT_API_URL: process.env.NEXT_PUBLIC_ROOT_API_URL,
  },
};

export default nextConfig;

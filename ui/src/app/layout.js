"use client"; // Add this line to indicate the file is a client component

import { Inter } from "next/font/google";
import "./globals.css";
import { useEffect } from "react";

export default function RootLayout({ children }) {
  useEffect(() => {
    document.documentElement.classList.add("dark");
  }, []);

  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}

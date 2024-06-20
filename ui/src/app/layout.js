"use client"; // Add this line to indicate the file is a client component

import "./globals.css";
import { useEffect } from "react";

import { Toaster } from "@/components/ui/sonner";

export default function RootLayout({ children }) {
  useEffect(() => {
    document.documentElement.classList.add("dark");
  }, []);

  return (
    <html lang="en">
      <head />
      <body>
        <main>{children}</main>
        <Toaster />
      </body>
    </html>
  );
}

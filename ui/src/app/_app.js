import { useEffect } from "react";
import "../styles/global.css"; // Make sure this path matches your file structure

function MyApp({ Component, pageProps }) {
  useEffect(() => {
    // Add 'dark' class to the HTML element to enable dark mode by default
    document.documentElement.classList.add("dark");
  }, []);

  return <Component {...pageProps} />;
}

export default MyApp;

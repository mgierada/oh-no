"use client"; // This is a client component 👈🏽
import { useEffect, useState } from "react";
import axios from "axios";

const CounterDisplay = () => {
  const [currentValue, setCurrentValue] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCurrentValue = async () => {
      try {
        // url = "https://ohno-server.fly.dev/counter"
        url = "http://localhost:3333/counter";
        const response = await axios.get(url);
        console.log("API response:", response);
        if (response.data && typeof response.data.CurrentValue === "number") {
          setCurrentValue(response.data.CurrentValue);
        } else {
          throw new Error("Invalid response data");
        }
      } catch (error) {
        console.error("Error fetching the current value:", error);
        setError("Failed to fetch current value");
      } finally {
        setLoading(false);
      }
    };

    fetchCurrentValue();
  }, []);

  return (
    <div>
      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
        Current Value: {currentValue}
      </h1>
      <div>
        {loading ? (
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
            Loading...
          </h1>
        ) : error ? (
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
            {error}
          </h1>
        ) : (
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
            Current streak: {currentValue} 🔥
          </h1>
        )}
      </div>
    </div>
  );
};

export default CounterDisplay;

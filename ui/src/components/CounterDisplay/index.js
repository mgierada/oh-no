import axios from "axios";

const fetchCurrentValue = async () => {
  try {
    const rootUrl = process.env.NEXT_PUBLIC_ROOT_API_URL;
    console.log("rootUrl:", rootUrl);
    const endpoint = `${rootUrl}/counter`;
    console.log("Endpoint:", endpoint);

    const response = await axios.get(endpoint);
    console.log("Response:", response.data); // Debug log
    if (response.data && typeof response.data.CurrentValue === "number") {
      return { currentValue: response.data.CurrentValue, error: null };
    } else {
      throw new Error("Invalid response data");
    }
  } catch (error) {
    console.error("Error fetching the current value:", error.message); // Debug log
    return { currentValue: null, error: "Failed to fetch current value" };
  }
};

const CounterDisplay = async () => {
  const { currentValue, error } = await fetchCurrentValue();

  return (
    <div>
      {error ? (
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          {error}
        </h1>
      ) : (
        <h1 className="scroll-m-20 text-4xl font-extrabold trackeng-tight lg:text-5xl">
          Current streak: {currentValue} 🔥
        </h1>
      )}
    </div>
  );
};

export default CounterDisplay;

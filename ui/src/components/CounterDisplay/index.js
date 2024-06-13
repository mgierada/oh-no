const CounterDisplay = ({ currentValue, error }) => {
  console.log("Rendered with currentValue:", currentValue, "and error:", error); // Debug log

  return (
    <div>
      {error ? (
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          {error}
        </h1>
      ) : (
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          Current streak: {currentValue} 🔥
        </h1>
      )}
    </div>
  );
};

export default CounterDisplay;

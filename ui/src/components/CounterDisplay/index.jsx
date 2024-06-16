"use client";

const CounterDisplay = ({ currentValue, error }) => {
  return (
    <div className="flex items-center justify-center h-full text-center">
      {error ? (
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          {error}
        </h1>
      ) : (
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          Current streak:
          <br /> {currentValue} days 🔥
        </h1>
      )}
    </div>
  );
};

export default CounterDisplay;

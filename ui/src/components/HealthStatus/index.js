"use client";

const HealthStatus = ({ isLocked, error }) => {
  let status = "";

  if (isLocked) {
    status = "Illness";
  } else {
    status = "Healthy";
  }

  return (
    <div className="flex items-center justify-center h-full text-center">
      {error ? (
        <h4 className="scroll-m-20 text-xl font-semibold tracking-tight">
          {error}
        </h4>
      ) : (
        <h4 className="scroll-m-20 text-xl font-semibold tracking-tight">
          {status}
        </h4>
      )}
    </div>
  );
};

export default HealthStatus;

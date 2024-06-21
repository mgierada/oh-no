import CounterDisplay from "@/components/CounterDisplay";
import Callendar from "@/components/Callendar";
import { DisplayCard } from "@/components/Card";
import HealthStatus from "@/components/HealthStatus";
import { ActionButton } from "@/components/ActionButton";
import { recordEvent } from "@/utils/actions";
import { ThemeToggle } from "@/components/ActionButton/ThemeToggle";

/**
 * Fetches the current value of the counter from the API.
 * @param {string} endpoint - The endpoint to fetch the counter data from.
 * @returns {Promise<Counter>} A promise that resolves to an object containing the current value and an error message if any.
 */
const fetchCounter = async (endpoint) => {
  try {
    const rootUrl = process.env.NEXT_PUBLIC_ROOT_API_URL;
    const url = `${rootUrl}/${endpoint}`;

    const response = await fetch(url, { cache: "no-store" });

    if (!response.ok) {
      throw new Error("Failed to fetch data");
    }

    const data = await response.json();

    if (!data) {
      throw new Error("Invalid response data");
    }

    const result = {
      currentValue: data.CurrentValue,
      maxValue: data.MaxValue,
      updatedAt: data.UpdatedAt,
      resetedAt: data.ResetedAt.String,
      wasEverReset: data.ResetedAt.Valid,
      isLocked: data.IsLocked,
      error: null,
    };
    return result;
  } catch (error) {
    return {
      currentValue: null,
      maxValue: null,
      updatedAt: null,
      resetedAt: null,
      wasEverReset: null,
      isLocked: false,
      error: `Failed to fetch ${endpoint} data.`,
    };
  }
};

/**
 * Home component.
 * @returns {JSX.Element} The Home component.
 */
const Home = async () => {
  const dataHealthyCounter = await fetchCounter("counter");
  const dataIllnessCounter = await fetchCounter("ohno-counter");

  const renderCounterDisplay = (dataHealthyCounter, dataIllnessCounter) => {
    if (dataHealthyCounter.error && dataIllnessCounter.error) {
      return (
        <CounterDisplay currentValue={null} error={dataHealthyCounter.error} />
      );
    }

    if (dataHealthyCounter.isLocked && dataIllnessCounter.isLocked) {
      return (
        <CounterDisplay
          currentValue="N/D"
          error="Something is wrong with counters."
        />
      );
    }

    if (!dataHealthyCounter.isLocked && !dataIllnessCounter.isLocked) {
      return (
        <CounterDisplay
          currentValue="N/D"
          error="Something is wrong with counters."
        />
      );
    }

    if (dataHealthyCounter.isLocked) {
      return (
        <CounterDisplay
          currentValue={dataIllnessCounter.currentValue}
          error={dataIllnessCounter.error}
        />
      );
    }

    return (
      <CounterDisplay
        currentValue={dataHealthyCounter.currentValue}
        error={dataHealthyCounter.error}
      />
    );
  };

  return (
    <main className="flex flex-col min-h-screen items-center justify-between p-4">
      <div className="flex space-x-4 justify-center lg:grid-cols-2 md:grid-cols-2">
        <DisplayCard
          title="Healthy"
          value={dataHealthyCounter.maxValue}
          description="Maximum number of days without illness"
        />
        <DisplayCard
          title="Illness"
          value={dataIllnessCounter.maxValue}
          description="Maximum duration of the illness in days"
        />
        <ThemeToggle />
      </div>
      <div className="flex flex-col items-center justify-center">
        <HealthStatus
          isLocked={dataHealthyCounter.isLocked}
          error={dataHealthyCounter.error}
        />
        {renderCounterDisplay(dataHealthyCounter, dataIllnessCounter)}
        <Callendar
          className="mt-5"
          lastTimeReseted={dataHealthyCounter.resetedAt}
          currentCouterValue={dataHealthyCounter.currentValue}
        />
        <p className="text-sm text-muted-foreground">
          Click dates to see the details.
        </p>
      </div>
      <div className="flex flex-row items-center justify-center gap-20">
        <ActionButton
          variant="destructive"
          toastMessage="Sick event successfully recorded"
          ctaButton="Sick"
          icon="biohazard"
          alertDialogDescription={`
            This action cannot be undone. 
            This will reset the healthy counter and start the sick interval.`}
          handleUpdate={async function update() {
            "use server";
            return recordEvent("ohno");
          }}
        />
        <ActionButton
          variant="outline"
          toastMessage="Recover event successfully recorded"
          ctaButton="Recover"
          icon="activity"
          alertDialogDescription={`
            This action cannot be undone. 
            This will reset the sick counter and start the healthy interval.`}
          handleUpdate={async function update() {
            "use server";
            return recordEvent("fine");
          }}
        />
      </div>
    </main>
  );
};

export default Home;

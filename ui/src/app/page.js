import CounterDisplay from "@/components/CounterDisplay";
import Callendar from "@/components/Callendar";
import { DisplayCard } from "@/components/Card";
import HealthStatus from "@/components/HealthStatus";
import { ActionButton } from "@/components/ActionButton";

/**
 * @typedef {Object} CounterApiResponse
 * @property {number} CurrentValue - The current value of the counter.
 * @property {number} MaxValue - the maximum value of the counter recorded so far.
 * @property {string} UpdatedAt - The timestamp when the counter was last updated.
 * @property {Object} ResetedAt - The reset information.
 * @property {string} ResetedAt.String - The timestamps when the counter was reset.
 * @property {boolean} ResetedAt.Valid - The validity of the reset object timestamps.
 * @property {boolean} IsLocked - Indicates if the counter is locked.
 */

/**
 * @typedef {Object} Counter
 * @property {number|null} currentValue - The current value of the counter.
 * @property {number|null} maxValue - The max value of the counter recorded so far.
 * @property {string|null} updatedAt - The timestamps when the last update was made.
 * @property {boolean} isLocked - A flag indicating whether the counter updates are currently locked
 * @property {string|null} resetedAt - The timestamps when the counter was reset.
 * @property {boolean|null} wasEverReset - The validity of the reset object timestamps.
 * @property {string|null} error - The error message if fetching failed.
 */

/**
 * @typedef {Object } RecordEventApiResponse
 * @property {string} message - The message returned from the API.
 */

/**
 * Records an event to the API.
 * @param {string} endpoint - The endpoint to record the event to.
 * @returns {Promise<void>} A promise that resolves when the event is successfully recorded.
 */
const recordEvent = async (endpoint) => {
  try {
    const rootUrl = process.env.NEXT_PUBLIC_ROOT_API_URL;
    const url = `${rootUrl}/${endpoint}`;

    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to record event");
    }

    /** @type {RecordEventApiResponse} */
    const data = await response.json();
    if (!data || !data.message) {
      throw new Error("Invalid response data");
    }
    return;
  } catch (error) {
    throw new Error(`Error recording event: ${error}`);
  }
};

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

    /** @type {CounterApiResponse} */
    const data = await response.json();

    if (!data) {
      throw new Error("Invalid response data");
    }

    /** @type {Counter} */
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
    /** @type {Counter} */
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
  /** @type {Counter} */
  const dataHealthyCounter = await fetchCounter("counter");
  /**@type {Counter} */
  const dataIllnessCounter = await fetchCounter("ohno-counter");

  /**
   * Method to render the counter display conditionally
   * @param {Counter} dataHealthyCounter - The healthy counter data.
   * @param {Counter} dataIllnessCounter - The illness counter data.
   * @returns {JSX.Element} The counter display component
   */
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
      </div>
      <div className="flex flex-col items-center justify-center">
        <HealthStatus
          isLocked={dataHealthyCounter.isLocked}
          error={dataHealthyCounter.error}
        />
        {renderCounterDisplay(dataHealthyCounter, dataIllnessCounter)}
        {/* <CounterDisplay */}
        {/*   currentValue={dataHealthyCounter.currentValue} */}
        {/*   error={dataHealthyCounter.error} */}
        {/* /> */}
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
        />
        <ActionButton
          variant="outline"
          toastMessage="Recover event successfully recorded"
          ctaButton="Recover"
          icon="activity"
          alertDialogDescription={`
            This action cannot be undone. 
            This will reset the sick counter and start the healthy interval.`}
        />
      </div>
    </main>
  );
};

export default Home;

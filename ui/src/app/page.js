import CounterDisplay from "@/components/CounterDisplay";
import Callendar from "@/components/Callendar";

/**
 * @typedef {Object} CounterApiResponse
 * @property {number} CurrentValue - The current value of the counter.
 * @property {string} UpdatedAt - The timestamp when the counter was last updated.
 * @property {Object} ResetedAt - The reset information.
 * @property {string} ResetedAt.String - The timestamps when the counter was reset.
 * @property {boolean} ResetedAt.Valid - The validity of the reset object timestamps.
 * @property {boolean} IsLocked - Indicates if the counter is locked.
 */

/**
 * @typedef {Object} Counter
 * @property {number|null} currentValue - The current value of the counter.
 * @property {string|null} updatedAt - The timestamps when the last update was made.
 * @property {boolean} isLocked - A flag indicating whether the counter updates are currently locked
 * @property {string|null} resetedAt - The timestamps when the counter was reset.
 * @property {boolean|null} wasEverReset - The validity of the reset object timestamps.
 * @property {string|null} error - The error message if fetching failed.
 */

/**
 * Fetches the current value of the counter from the API.
 * @returns {Promise<Counter>} A promise that resolves to an object containing the current value and an error message if any.
 */
const fetchCounter = async () => {
  try {
    const rootUrl = process.env.NEXT_PUBLIC_ROOT_API_URL;
    const endpoint = `${rootUrl}/counter`;

    const response = await fetch(endpoint, { cache: "no-store" });

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
      updatedAt: null,
      resetedAt: null,
      wasEverReset: null,
      isLocked: false,
      error: "Failed to fetch current value",
    };
  }
};

/**
 * Home component.
 * @returns {JSX.Element} The Home component.
 */
const Home = async () => {
  /** @type {Counter} */
  const data = await fetchCounter();

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-5">
      <div className="z-10 max-w-5xl w-full flex flex-col items-center justify-center font-mono text-sm lg:flex"></div>
      <div className="flex flex-col items-center justify-center">
        <CounterDisplay currentValue={data.currentValue} error={data.error} />
      </div>
      <div className="flex flex-col items-center justify-center">
        <Callendar
          className="mt-5"
          lastTimeReseted={data.resetedAt}
          currentCouterValue={data.currentValue}
        />
      </div>
    </main>
  );
};

export default Home;

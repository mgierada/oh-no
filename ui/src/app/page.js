import axios from "axios";
import CounterDisplay from "@/components/CounterDisplay";

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

    const response = await axios.get(endpoint);

    /** @type {CounterApiResponse} */
    const data = response.data;
    console.log("Response data:", data);

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
    console.error("Error fetching the current value:", error.message);
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
  const { currentValue, error } = await fetchCounter();

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex"></div>

      <div className="relative flex place-items-center before:absolute before:h-[300px] before:w-full sm:before:w-[480px] before:-translate-x-1/2 before:rounded-full before:bg-gradient-radial before:from-white before:to-transparent before:blur-2xl before:content-[''] after:absolute after:-z-20 after:h-[180px] after:w-full sm:after:w-[240px] after:translate-x-1/3 after:bg-gradient-conic after:from-sky-200 after:via-blue-200 after:blur-2xl after:content-[''] before:dark:bg-gradient-to-br before:dark:from-transparent before:dark:to-blue-700 before:dark:opacity-10 after:dark:from-sky-900 after:dark:via-[#0141ff] after:dark:opacity-40 before:lg:h-[360px] z-[-1]">
        <CounterDisplay currentValue={currentValue} error={error} />
      </div>
      <div></div>
    </main>
  );
};

export default Home;

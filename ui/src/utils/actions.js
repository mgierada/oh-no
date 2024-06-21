"use server";

/**
 * @typedef {Object } RecordEventApiResponse
 * @property {string} message - The message returned from the API.
 */

/**
 * Records an event to the API.
 * @param {string} endpoint - The endpoint to record the event to.
 * @returns {Promise<void>} A promise that resolves when the event is successfully recorded.
 */
export const recordEvent = async (endpoint) => {
  // "use server";
  try {
    const rootUrl = process.env.NEXT_PUBLIC_ROOT_API_URL;
    const url = `${rootUrl}/${endpoint}`;

    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    console.log("response: ", response);

    if (!response.ok) {
      throw new Error("Failed to record event");
    }

    const data = await response.json();
    if (!data || !data.message) {
      throw new Error("Invalid response data");
    }
    return;
  } catch (error) {
    console.log(error);
    throw new Error(`Error recording event: ${error.message}`);
  }
};

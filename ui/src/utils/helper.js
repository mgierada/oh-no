/*
 * This function formats the current date and time in a human-readable format.
 * @returns {string} formattedNow - The formatted date and time in the following format
 * Thursday, June 20 at 2024 at 2:41 PM
 */
const formatCurrentDate = () => {
  const options = {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "2-digit",
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  };
  const now = new Date();
  return now.toLocaleString("en-US", options).replace(/,([^,]*)$/, " at$1");
};

export default formatCurrentDate;

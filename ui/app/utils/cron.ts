import { CronExpressionParser } from "cron-parser";

function getTime(cron: string): string {
  try {
    const interval = CronExpressionParser.parse(cron);
    const nextDate = interval.next().toDate();
    const time = nextDate.toLocaleTimeString("en-SG", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
    return time;
  } catch (err) {
    console.error("Invalid cron expression", err);
    return "00:00";
  }
}

function getFrequency(
  cron: string,
): "custom" | "daily" | "weekdays" | "weekends" {
  try {
    const parts = cron.trim().split(" ");

    if (parts.length < 5) return "custom";

    if (parts[4] === "*") return "daily";
    if (parts[4] === "1-5") return "weekdays";
    if (parts[4] === "0,6" || parts[4] === "6,0") return "weekends";

    return "custom";
  } catch (err) {
    console.error("Invalid cron expression", err);
    return "custom";
  }
}

export { getTime, getFrequency };

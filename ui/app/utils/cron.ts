import { CronExpressionParser } from "cron-parser";

function getTime(cron: string): string {
  try {
    const interval = CronExpressionParser.parse(cron, { tz: "UTC" });
    const time = interval.next().toDate().toLocaleTimeString("en-SG", {
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

function getFrequency(cron: string): string {
  try {
    const parts = cron.trim().split(" ");
    if (parts.length < 5) return "-";
    return parts[4] as string;
  } catch (err) {
    console.error("Invalid cron expression", err);
    return "-";
  }
}

function getCron(time: string, frequency: string): string {
  const hour: number = time.split(":").map(Number)[0] ?? 0;

  const date = new Date();
  date.setHours(hour, 0, 0, 0);

  return `0 ${date.getUTCHours()} * * ${frequency}`;
}

export { getTime, getFrequency, getCron };

import { init, initData } from "@telegram-apps/sdk-vue";

const checkTelegramEnvironment = (): boolean => {
  // Check for Telegram WebApp API
  const hasTelegramWebApp = !!(window as any)?.Telegram?.WebApp;

  // Check for Telegram-specific URL parameters
  const urlParams = new URLSearchParams(window.location.search);
  const hasTelegramParams =
    urlParams.has("tgWebAppPlatform") ||
    urlParams.has("tgWebAppVersion") ||
    window.location.hash.includes("tgWebApp");

  // Check user agent for Telegram
  const userAgent = navigator.userAgent.toLowerCase();
  const isTelegramUserAgent = userAgent.includes("telegram");

  return hasTelegramWebApp || hasTelegramParams || isTelegramUserAgent;
};

const getChatId = (): number | undefined => {
  init();
  initData.restore();
  return initData.chat()?.id;
};

const getInitDataRaw = (): string | undefined => {
  init();
  initData.restore();
  return initData.raw();
};

const getTheme = (): "light" | "dark" => {
  // Check for saved theme preference
  const savedTheme = localStorage.getItem("alfred-theme") as "light" | "dark";
  if (savedTheme) {
    return savedTheme;
  } else {
    // Check system preference
    const prefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)"
    ).matches;
    return prefersDark ? "dark" : "light";
  }
};

const setTheme = (theme: "light" | "dark"): void => {
  document.documentElement.classList.remove("light", "dark");
  document.documentElement.classList.add(theme);
  localStorage.setItem("alfred-theme", theme);
};

export {
  checkTelegramEnvironment,
  getChatId,
  getInitDataRaw,
  getTheme,
  setTheme,
};

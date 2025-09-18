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

export { checkTelegramEnvironment };

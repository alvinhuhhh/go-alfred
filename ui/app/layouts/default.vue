<script setup lang="ts">
const appConfig = useAppConfig();
const chatId = useState<number>("chatId", () => {
  if (appConfig.devMode) {
    return 1;
  }
  return getChatId() ?? getUserId() ?? 1;
});
const initDataRaw = useState<string>("initDataRaw", () => {
  if (appConfig.devMode) {
    return "";
  }
  return getInitDataRaw() ?? "";
});
const encryptionKey = useState<string>("encryptionKey", () => {
  return "";
});

onMounted(async () => {
  setTheme(getTheme());
  encryptionKey.value = appConfig.devMode
    ? "+tuHPldsy0hy16yebxkQsmlHiZKkhlq3gzm447tWdkQ="
    : await fetchEncryptionKey(
        appConfig.keyVersion,
        chatId.value,
        initDataRaw.value
      );
});
</script>

<template>
  <ToastProvider>
    <slot />
    <ToastViewport />
  </ToastProvider>
</template>

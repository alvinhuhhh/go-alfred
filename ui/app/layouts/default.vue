<script setup lang="ts">
const { public: config } = useRuntimeConfig();
const chatId = useState<number>("chatId", () => {
  if (import.meta.dev) {
    return 1;
  }
  return getChatId() ?? getUserId() ?? 1;
});
const initDataRaw = useState<string>("initDataRaw", () => {
  if (import.meta.dev) {
    return "";
  }
  return getInitDataRaw() ?? "";
});
const encryptionKey = useState<string>("encryptionKey", () => {
  return "";
});

onMounted(async () => {
  setTheme(getTheme());
  encryptionKey.value = import.meta.dev
    ? "+tuHPldsy0hy16yebxkQsmlHiZKkhlq3gzm447tWdkQ="
    : await fetchEncryptionKey(
        config.keyVersion as number,
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

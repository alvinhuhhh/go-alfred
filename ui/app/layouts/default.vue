<script setup lang="ts">
const runtimeConfig = useRuntimeConfig();
const chatId = useState<number>("chatId", () => {
  if (runtimeConfig.public.devMode) {
    return 1;
  }
  return getChatId() ?? getUserId() ?? 1;
});
const initDataRaw = useState<string>("initDataRaw", () => {
  if (runtimeConfig.public.devMode) {
    return "";
  }
  return getInitDataRaw() ?? "";
});
const encryptionKey = useState<string>("encryptionKey", () => {
  return "";
});

onMounted(async () => {
  console.log(runtimeConfig.public.devMode);
  console.log(runtimeConfig.public.keyVersion);
  setTheme(getTheme());
  encryptionKey.value = runtimeConfig.public.devMode
    ? "+tuHPldsy0hy16yebxkQsmlHiZKkhlq3gzm447tWdkQ="
    : await fetchEncryptionKey(
        runtimeConfig.keyVersion as number,
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

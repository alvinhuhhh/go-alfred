<script setup lang="ts">
const route = useRoute();
const chatId = useState<number>("chatId", () => {
  if (import.meta.env.VITE_ENV === "production") {
    if (route.path.includes("/telegram")) {
      return getChatId() ?? getUserId() ?? 1;
    }
  }
  return 1;
});
const initDataRaw = useState<string>("initDataRaw", () => {
  if (import.meta.env.VITE_ENV === "production") {
    if (route.path.includes("/telegram")) {
      return getInitDataRaw() ?? "";
    }
  }
  return "";
});
const encryptionKey = useState<string>("encryptionKey", () => {
  return "";
});

onMounted(async () => {
  setTheme(getTheme());
  if (route.path.includes("/telegram")) {
    encryptionKey.value = await fetchEncryptionKey(
      import.meta.env.VITE_MASTER_KEY_VERSION as number,
      chatId.value,
      initDataRaw.value
    );
  }
});
</script>

<template>
  <ToastProvider>
    <slot />
    <ToastViewport />
  </ToastProvider>
</template>

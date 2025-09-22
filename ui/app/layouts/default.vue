<script setup>
import { ToastProvider } from "reka-ui";

onMounted(() => {
  setTheme(getTheme());
});
</script>

<template>
  <div>
    <ToastProvider>
      <slot />
    </ToastProvider>
  </div>
</template>

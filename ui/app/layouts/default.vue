<script setup>
onMounted(() => {
  setTheme(getTheme());
});
</script>

<template>
  <ToastProvider>
    <slot />
    <ToastViewport />
  </ToastProvider>
</template>

<script setup>
import { ToastProvider } from "reka-ui";
</script>

<template>
  <ToastProvider>
    <slot />
  </ToastProvider>
</template>

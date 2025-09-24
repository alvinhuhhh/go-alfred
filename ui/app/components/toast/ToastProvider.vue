<script setup>
import { ToastProvider } from "reka-ui";
</script>

<template>
  <ToastProvider swipe-direction="down">
    <slot />
  </ToastProvider>
</template>

<script setup>
import { ToastProvider } from "reka-ui";
</script>

<template>
  <div>
    <ToastProvider>
      <slot />
    </ToastProvider>
  </div>
</template>

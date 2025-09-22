<script setup>
import { DialogOverlay } from "reka-ui";
</script>

<template>
  <DialogOverlay
    class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/50"
  >
    <slot />
  </DialogOverlay>
</template>

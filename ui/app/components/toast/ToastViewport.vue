<script setup>
import { ToastViewport } from "reka-ui";
</script>

<template>
  <ToastViewport
    class="[--viewport-padding:_25px] fixed inset-0 m-auto flex flex-col p-[var(--viewport-padding)] gap-[10px] w-[390px] max-w-[100vw] list-none z-[2147483647] outline-none"
  >
    <slot />
  </ToastViewport>
</template>

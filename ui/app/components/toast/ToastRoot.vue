<script setup>
import { XIcon } from "lucide-vue-next";
import { ToastRoot } from "reka-ui";
</script>

<template>
  <ToastRoot
    class="bg-background rounded-lg shadow-sm border p-3 flex justify-between items-center data-[state=open]:animate-slideUp data-[state=closed]:animate-slideDown"
  >
    <slot />
    <ToastClose>
      <Button variant="ghost" size="sm">
        <XIcon class="w-4 h-4 text-foreground" />
      </Button>
    </ToastClose>
  </ToastRoot>
</template>

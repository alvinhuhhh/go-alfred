<script setup>
import { ToastDescription } from "reka-ui";
import { BadgeCheckIcon, CircleSlashIcon, InfoIcon } from "lucide-vue-next";

const props = defineProps({
  status: "success" | "error" | "info",
});

const icon = computed(() => {
  switch (props.status) {
    case "success":
      return BadgeCheckIcon;
    case "error":
      return CircleSlashIcon;
    case "info":
      return InfoIcon;
    default:
      return BadgeCheckIcon;
  }
});
</script>

<template>
  <ToastDescription class="flex items-center gap-x-2">
    <component :is="icon" class="w-4 h-4 ml-4 text-foreground" />
    <span class="text-sm text-foreground">
      <slot />
    </span>
  </ToastDescription>
</template>

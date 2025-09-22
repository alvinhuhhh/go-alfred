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
  <ToastDescription
    class="[grid-area:_description] m-0 text-slate11 text-xs leading-[1.3] flex items-center gap-2"
  >
    <component :is="icon" class="w-4 h-4" />
    <span class="text-sm font-bold text-muted-foreground">
      <slot />
    </span>
  </ToastDescription>
</template>

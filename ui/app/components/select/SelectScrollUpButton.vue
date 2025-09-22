<script setup>
import { SelectScrollUpButton } from "reka-ui";
import { ChevronDownIcon } from "lucide-vue-next";
</script>

<template>
  <SelectScrollUpButton
    class="flex cursor-default items-center justify-center py-1"
  >
    <ChevronDownIcon class="size-4" />
    <slot />
  </SelectScrollUpButton>
</template>

<script setup>
import { SelectScrollDownButton } from "reka-ui";
import { ChevronUpIcon } from "lucide-vue-next";
</script>

<template>
  <SelectScrollDownButton
    class="flex cursor-default items-center justify-center py-1"
  >
    <ChevronUpIcon class="size-4" />
    <slot />
  </SelectScrollDownButton>
</template>

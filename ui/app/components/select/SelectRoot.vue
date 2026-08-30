<script setup>
import { SelectRoot } from "reka-ui";
const model = defineModel();
</script>

<template>
  <SelectRoot v-model="model">
    <slot />
  </SelectRoot>
</template>

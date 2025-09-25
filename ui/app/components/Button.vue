<script setup>
const props = defineProps({
  variant:
    "default" | "secondary" | "destructive" | "outline" | "ghost" | "link",
  size: "default" | "sm" | "lg" | "icon",
});

const variant = computed(() => {
  switch (props.variant) {
    case "destructive":
      return "bg-destructive text-white hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60";
    case "outline":
      return "border bg-background text-foreground hover:bg-accent-hover hover:text-accent-foreground dark:bg-input dark:border-border dark:hover:bg-accent-hover";
    case "secondary":
      return "bg-secondary text-secondary-foreground hover:bg-secondary-hover";
    case "ghost":
      return "hover:bg-accent-hover hover:text-accent-foreground";
    case "link":
      return "text-primary underline-offset-4 hover:underline";
    default:
      return "bg-primary text-primary-foreground hover:bg-primary/90";
  }
});

const size = computed(() => {
  switch (props.size) {
    case "sm":
      return "h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5";
    case "lg":
      return "h-10 rounded-md px-6 has-[>svg]:px-4";
    case "icon":
      return "size-9 rounded-md";
    default:
      return "h-9 px-4 py-2 has-[>svg]:px-3";
  }
});
</script>

<template>
  <button
    class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive"
    :class="[variant, size]"
  >
    <slot />
  </button>
</template>

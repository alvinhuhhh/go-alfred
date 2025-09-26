<script setup>
import { User, BookOpen, Settings, Sparkles } from "lucide-vue-next";

const appConfig = useAppConfig();
if (!appConfig.devMode && !checkTelegramEnvironment()) {
  // If not opened in Telegram redirect to Landing
  navigateTo("/");
}

const chatId = 1;

const features = [
  {
    id: "notebook",
    title: "Notebook",
    description: "Manage shared secrets and notes",
    icon: BookOpen,
    color: "bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400",
    available: true,
  },
  {
    id: "settings",
    title: "Settings",
    description: "Configure Alfred preferences",
    icon: Settings,
    color: "bg-gray-50 dark:bg-gray-950 text-gray-600 dark:text-gray-400",
    available: false,
  },
  {
    id: "more",
    title: "More Features",
    description: "More tools coming soon",
    icon: Sparkles,
    color:
      "bg-purple-50 dark:bg-purple-950 text-purple-600 dark:text-purple-400",
    available: false,
  },
];

function handleFeatureClick(id) {
  switch (id) {
    case "notebook":
      return navigateTo("/telegram/notebook/" + chatId);
    case "settings":
      return navigateTo("/telegram/settings/" + chatId);
    default:
      return;
  }
}
</script>

<template>
  <div class="min-h-screen bg-background p-4">
    <div class="max-w-md mx-auto">
      <!-- Header -->
      <div class="text-center mb-8 pt-4">
        <div
          class="w-24 h-24 bg-secondary border rounded-full mx-auto mb-2 flex items-center justify-center"
        >
          <img
            src="~/assets/img/alfred-v2-transparent-bg.png"
            class="w-full h-full rounded-full"
          />
        </div>
        <h1 class="text-2xl font-medium text-foreground mb-2">Alfred</h1>
        <p class="text-muted-foreground">
          Your shared assistant powered by Telegram
        </p>
      </div>

      <!-- Feature Cards -->
      <div class="space-y-4">
        <Card
          v-for="feature in features"
          :key="feature.id"
          class="p-0 overflow-hidden"
        >
          <Button
            variant="ghost"
            class="w-full h-auto p-6 flex items-center justify-start space-x-4 hover:bg-accent-hover"
            :class="
              feature.available
                ? ''
                : 'opacity-60 cursor-not-allowed hover:bg-transparent'
            "
            @click="handleFeatureClick(feature.id)"
            :disabled="!feature.available"
          >
            <div
              class="w-12 h-12 rounded-lg flex items-center justify-center"
              :class="feature.color"
            >
              <component :is="feature.icon" class="w-6 h-6" />
            </div>
            <div class="flex-1 text-left">
              <div class="flex items-center space-x-2 mb-1">
                <h3 class="font-medium text-foreground">
                  {{ feature.title }}
                </h3>
                <Badge
                  v-if="!feature.available"
                  variant="secondary"
                  class="text-xs px-2 py-0.5"
                >
                  Coming Soon
                </Badge>
              </div>
              <p class="text-sm text-muted-foreground">
                {{ feature.description }}
              </p>
            </div>
          </Button>
        </Card>
      </div>

      <!-- Footer -->
      <div class="mt-12 text-center">
        <p class="text-xs text-muted-foreground">
          Alfred v1.0 - Your digital butler
        </p>
      </div>
    </div>
  </div>
</template>

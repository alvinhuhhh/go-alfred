<script setup>
import {
  MessageSquare,
  ArrowRight,
  Smartphone,
  BookOpen,
  Sparkles,
} from "lucide-vue-next";

const runtimeConfig = useRuntimeConfig();
if (!runtimeConfig.public.devMode && checkTelegramEnvironment()) {
  // Redirect to /telegram if opened in Telegram
  navigateTo("/telegram");
}

const features = [
  {
    title: "Secure Notebook",
    description:
      "Store and share passwords, notes, and important information securely with your group",
    icon: BookOpen,
    color: "bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400",
  },
  {
    title: "Smart Scheduling",
    description:
      "Automated reminders and messages sent to your group chat at perfect times",
    icon: MessageSquare,
    color: "bg-green-50 dark:bg-green-950 text-green-600 dark:text-green-400",
  },
  {
    title: "Smart Automation",
    description:
      "Coming soon: Advanced tools for managing shared tasks and reminders",
    icon: Sparkles,
    color:
      "bg-purple-50 dark:bg-purple-950 text-purple-600 dark:text-purple-400",
  },
];

function handleGetStarted() {
  console.log("get-started");
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <div class="max-w-4xl mx-auto">
      <!-- Hero Section -->
      <div class="max-w-4xl mx-auto px-4 py-16">
        <div class="text-center">
          <div
            class="w-24 h-24 bg-secondary border rounded-full mx-auto mb-4 flex items-center justify-center"
          >
            <img
              src="~/assets/img/alfred-v2-transparent-bg.png"
              class="w-full h-full rounded-full"
            />
          </div>
          <h1 class="text-4xl md:text-6xl font-medium text-foreground mb-4">
            Meet Alfred
          </h1>
          <p class="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            Your personal administrative assistant. Organize, automate, and
            streamline your group's shared tasks through Telegram.
          </p>
          <div
            class="flex flex-col sm:flex-row gap-4 justify-center items-center"
          >
            <Button
              @click="handleGetStarted"
              size="lg"
              class="flex items-center space-x-2"
            >
              <MessageSquare class="w-5 h-5" />
              <span>Start with Telegram</span>
              <ArrowRight class="w-4 h-4" />
            </Button>
            <Badge variant="secondary" class="flex items-center space-x-1">
              <Smartphone class="w-3 h-3" />
              <span>Telegram Mini App</span>
            </Badge>
          </div>
        </div>
      </div>

      <!-- Features Grid -->
      <div class="grid md:grid-cols-3 gap-8 mb-16">
        <Card
          v-for="feature in features"
          :key="feature.title"
          class="p-6 text-center"
        >
          <div
            class="w-16 h-16 rounded-lg flex items-center justify-center mx-auto mb-4"
            :class="feature.color"
          >
            <component :is="feature.icon" />
          </div>
          <h3 class="text-xl font-medium text-foreground mb-3">
            {{ feature.title }}
          </h3>
          <p class="text-muted-foreground">{{ feature.description }}</p>
        </Card>
      </div>

      <!-- How it works -->
      <div class="text-center mb-16">
        <h2 class="text-3xl font-medium text-foreground mb-8">
          How Alfred Works
        </h2>
        <div class="grid md:grid-cols-3 gap-8">
          <div class="space-y-4">
            <div
              class="w-12 h-12 bg-blue-50 dark:bg-blue-950 rounded-full flex items-center justify-center mx-auto"
            >
              <span class="text-blue-600 dark:text-blue-400 font-medium"
                >1</span
              >
            </div>
            <h4 class="font-medium text-foreground">Add to Telegram</h4>
            <p class="text-muted-foreground">
              Start a conversation with Alfred on Telegram
            </p>
          </div>
          <div class="space-y-4">
            <div
              class="w-12 h-12 bg-green-50 dark:bg-green-950 rounded-full flex items-center justify-center mx-auto"
            >
              <span class="text-green-600 dark:text-green-400 font-medium"
                >2</span
              >
            </div>
            <h4 class="font-medium text-foreground">Configure Settings</h4>
            <p class="text-muted-foreground">
              Set up your preferences and group settings
            </p>
          </div>
          <div class="space-y-4">
            <div
              class="w-12 h-12 bg-purple-50 dark:bg-purple-950 rounded-full flex items-center justify-center mx-auto"
            >
              <span class="text-purple-600 dark:text-purple-400 font-medium"
                >3</span
              >
            </div>
            <h4 class="font-medium text-foreground">Enjoy Automation</h4>
            <p class="text-muted-foreground">
              Let Alfred handle your group's administrative tasks
            </p>
          </div>
        </div>
      </div>

      <!-- CTA Section -->
      <Card class="p-8 text-center bg-muted/50">
        <h3 class="text-2xl font-medium text-foreground mb-4">
          Ready to Get Started?
        </h3>
        <p class="text-muted-foreground mb-6">
          Your simple Telegram assistant for dinner plans and shared notes.
        </p>
        <Button
          @click="handleGetStarted"
          size="lg"
          class="flex items-center space-x-2 mx-auto"
        >
          <MessageSquare class="w-5 h-5" />
          <span>Open in Telegram</span>
          <ArrowRight class="w-4 h-4" />
        </Button>
      </Card>

      <!-- Footer -->
      <div class="text-center py-16 text-muted-foreground">
        <p>Alfred v1.0 - Your digital butler, powered by Telegram</p>
      </div>
    </div>
  </div>
</template>

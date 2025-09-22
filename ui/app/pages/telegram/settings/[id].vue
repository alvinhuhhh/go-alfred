<script setup>
import { ArrowLeft, Moon, Sun, MessageSquare, Clock } from "lucide-vue-next";

const theme = "light";
const scheduleEnabled = ref(true);
const scheduleTime = ref(null);
const scheduleFrequency = ref(null);

const timeOptions = [
  { value: "06:00", label: "6:00 AM" },
  { value: "07:00", label: "7:00 AM" },
  { value: "08:00", label: "8:00 AM" },
  { value: "09:00", label: "9:00 AM" },
  { value: "10:00", label: "10:00 AM" },
  { value: "11:00", label: "11:00 AM" },
  { value: "12:00", label: "12:00 PM" },
  { value: "13:00", label: "1:00 PM" },
  { value: "14:00", label: "2:00 PM" },
  { value: "15:00", label: "3:00 PM" },
  { value: "16:00", label: "4:00 PM" },
  { value: "17:00", label: "5:00 PM" },
  { value: "18:00", label: "6:00 PM" },
  { value: "19:00", label: "7:00 PM" },
  { value: "20:00", label: "8:00 PM" },
  { value: "21:00", label: "9:00 PM" },
  { value: "22:00", label: "10:00 PM" },
];

const frequencyOptions = [
  { value: "daily", label: "Daily" },
  { value: "weekdays", label: "Weekdays Only" },
  { value: "weekends", label: "Weekends Only" },
  { value: "monday", label: "Every Monday" },
  { value: "tuesday", label: "Every Tuesday" },
  { value: "wednesday", label: "Every Wednesday" },
  { value: "thursday", label: "Every Thursday" },
  { value: "friday", label: "Every Friday" },
  { value: "saturday", label: "Every Saturday" },
  { value: "sunday", label: "Every Sunday" },
];

function back() {
  return navigateTo("/telegram");
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Header -->
    <div class="sticky top-0 bg-background border-b border-border p-4">
      <div class="flex items-center justify-between max-w-md mx-auto">
        <div class="flex items-center space-x-3">
          <Button @click="back" variant="ghost" size="sm" class="p-2">
            <ArrowLeft class="w-5 h-5" />
          </Button>
          <h1 class="text-xl font-medium">Settings</h1>
        </div>
      </div>
    </div>

    <!-- Settings Content -->
    <div class="p-4">
      <div class="max-w-md mx-auto space-y-6">
        <!-- Scheduled Messages Section -->
        <Card class="p-6">
          <div class="flex items-center space-x-3 mb-4">
            <div
              class="w-10 h-10 bg-blue-50 dark:bg-blue-950 rounded-lg flex items-center justify-center"
            >
              <MessageSquare class="w-5 h-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h3 class="font-medium">Scheduled Messages</h3>
              <p class="text-sm text-muted-foreground">
                Send automatic reminders to your family chat
              </p>
            </div>
          </div>

          <div class="space-y-6">
            <!-- Enable/Disable Toggle -->
            <div class="flex items-center justify-between">
              <div>
                <Label htmlFor="schedule-enabled" class="text-base"
                  >Enable Schedule</Label
                >
                <p class="text-sm text-muted-foreground">
                  Turn scheduled messages on or off
                </p>
              </div>
              <Switch
                id="schedule-enabled"
                checked="{scheduleEnabled}"
                onCheckedChange="{setScheduleEnabled}"
              />
            </div>

            <!-- Time Selection -->
            <div
              class="space-y-2"
              :class="!scheduleEnabled ? 'opacity-50' : ''"
            >
              <Label for="schedule-time">Time</Label>
              <SelectRoot :value="scheduleTime" :disabled="!scheduleEnabled">
                <SelectTrigger id="schedule-time">
                  <div class="flex items-center space-x-2">
                    <Clock class="w-4 h-4" />
                    <SelectValue placeholder="Select time" />
                  </div>
                </SelectTrigger>
                <SelectContent>
                  <SelectViewport>
                    <SelectItem
                      v-for="option in timeOptions"
                      :key="option.value"
                      :value="option.value"
                    >
                      <SelectItemText>
                        {{ option.label }}
                      </SelectItemText>
                    </SelectItem>
                  </SelectViewport>
                </SelectContent>
              </SelectRoot>
            </div>

            <!-- Frequency Selection -->
            <div
              class="space-y-2"
              :class="!scheduleEnabled ? 'opacity-50' : ''"
            >
              <Label for="schedule-frequency">Frequency</Label>
              <SelectRoot
                :value="scheduleFrequency"
                :disabled="!scheduleEnabled"
              >
                <SelectTrigger id="schedule-frequency">
                  <SelectValue placeholder="Select frequency" />
                </SelectTrigger>
                <SelectContent>
                  <SelectViewport>
                    <SelectItem
                      v-for="option in frequencyOptions"
                      :key="option.value"
                      :value="option.value"
                    >
                      <SelectItemText>
                        {{ option.label }}
                      </SelectItemText>
                    </SelectItem>
                  </SelectViewport>
                </SelectContent>
              </SelectRoot>
            </div>

            <!-- Current Schedule Summary -->
            <div class="p-3 bg-muted rounded-lg">
              <p class="text-sm">
                <span class="font-medium">Current schedule:</span> Messages will
                be sent
                <span class="font-medium"></span>
                <span class="font-medium"> </span>
              </p>
            </div>
          </div>
        </Card>

        <!-- App Appearance Section -->
        <Card class="p-6">
          <div class="flex items-center space-x-3 mb-4">
            <div
              class="w-10 h-10 bg-purple-50 dark:bg-purple-950 rounded-lg flex items-center justify-center"
            >
              <Moon
                v-if="theme === 'dark'"
                class="w-5 h-5 text-purple-600 dark:text-purple-400"
              />
              <Sun
                v-else
                class="w-5 h-5 text-purple-600 dark:text-purple-400"
              />
            </div>
            <div>
              <h3 class="font-medium">Appearance</h3>
              <p class="text-sm text-muted-foreground">
                Customize the app's visual theme
              </p>
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div>
              <Label htmlFor="dark-mode" class="text-base">Dark Mode</Label>
              <p class="text-sm text-muted-foreground">
                Switch between light and dark themes
              </p>
            </div>
            <Switch id="dark-mode" />
          </div>
        </Card>

        <!-- Save Settings -->
        <Button class="w-full"> Save Settings </Button>

        <!-- App Info Section -->
        <Card class="p-6 bg-transparent">
          <h3 class="font-medium mb-3">About Alfred</h3>
          <div class="space-y-2 text-sm text-muted-foreground">
            <p>Version 1.0.0</p>
            <p>Built for family organization and management</p>
            <p>Telegram Mini App</p>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

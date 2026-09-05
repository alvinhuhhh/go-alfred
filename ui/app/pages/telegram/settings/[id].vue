<script setup lang="ts">
import { ArrowLeft, Moon, Sun, MessageSquare, Clock } from "lucide-vue-next";
import { getCron, getFrequency, getSummary, getTime } from "~/utils/cron";

interface Schedule {
  enabled: boolean;
  time: string;
  frequency: string;
  summary: string | null;
}

const { public: config } = useRuntimeConfig();

const route = useRoute();
const chatId: string = route.params.id as string;
const initDataRaw = useState<string>("initDataRaw");
const defaultSchedule: Schedule = {
  enabled: false,
  time: "00:00",
  frequency: "daily",
  summary: "",
};

const initialSchedule = ref<Schedule>(defaultSchedule);
const schedule = ref<Schedule>(defaultSchedule);
const isDarkMode = ref(getTheme() === "dark");
const isChanged = computed<boolean>(() => {
  return (
    JSON.stringify(schedule.value) !== JSON.stringify(initialSchedule.value)
  );
});

const {
  data: raw,
  pending,
  error,
  refresh,
} = await useAsyncData<string>("cron", () =>
  $fetch(`/api/settings/cron/dinner-${chatId}`, {
    method: "GET",
    headers: {
      Authorization: `tma ${initDataRaw.value}`,
    },
  }),
);

watch(
  raw,
  (raw) => {
    if (!raw) return;
    const json = JSON.parse(raw);
    const time = getTime(json.cron);
    const frequency = getFrequency(json.cron);

    const data: Schedule = {
      enabled: true,
      time: time,
      frequency: frequency,
      summary: getSummary(time, frequency),
    };
    initialSchedule.value = structuredClone(data);
    schedule.value = data;
  },
  { immediate: true },
);

const timeOptions = [
  { value: "00:00", label: "0:00 AM" },
  { value: "01:00", label: "1:00 AM" },
  { value: "02:00", label: "2:00 AM" },
  { value: "03:00", label: "3:00 AM" },
  { value: "04:00", label: "4:00 AM" },
  { value: "05:00", label: "5:00 AM" },
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
  { value: "23:00", label: "11:00 PM" },
];

const frequencyOptions = [
  { value: "*", label: "Daily" },
  { value: "1-5", label: "Weekdays Only" },
  { value: "0,6", label: "Weekends Only" },
];

function back() {
  return navigateTo("/telegram");
}

function toggleTheme() {
  setTheme(getTheme() === "light" ? "dark" : "light");
}

async function saveSettings() {
  console.debug(schedule.value);
  console.debug(initialSchedule.value);
  if (!isChanged.value) return navigateTo("/telegram");

  if (!schedule.value.enabled) {
    console.debug("Schedule disabled, sending DELETE request");
    await $fetch(`/api/settings/cron/dinner-${chatId}`, {
      method: "DELETE",
      headers: {
        Authorization: `tma ${initDataRaw.value}`,
      },
    });
  } else {
    console.debug("Schedule updated, sending POST request");
    const cron = getCron(schedule.value.time, schedule.value.frequency);
    const body = {
      url: useRequestURL().hostname + "/api/cron",
      chatId: parseInt(chatId),
      jobName: `dinner-${chatId}`,
      schedule: cron,
    };

    await $fetch("/api/settings/cron", {
      method: "POST",
      body: JSON.stringify(body),
      headers: {
        Authorization: `tma ${initDataRaw.value}`,
      },
    });
  }

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
                <Label for="schedule-enabled" class="text-base"
                  >Enable Schedule</Label
                >
                <p class="text-sm text-muted-foreground">
                  Turn scheduled messages on or off
                </p>
              </div>
              <Switch id="schedule-enabled" v-model="schedule.enabled" />
            </div>

            <!-- Time Selection -->
            <div
              class="space-y-2"
              :class="!schedule.enabled ? 'opacity-50' : ''"
            >
              <Label for="schedule-time">Time</Label>
              <SelectRoot v-model="schedule.time" :disabled="!schedule.enabled">
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
              :class="!schedule.enabled ? 'opacity-50' : ''"
            >
              <Label for="schedule-frequency">Frequency</Label>
              <SelectRoot
                v-model="schedule.frequency"
                :disabled="!schedule.enabled"
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
            <div class="p-3 bg-muted rounded-lg" :hidden="!schedule.enabled">
              <p class="text-sm">
                <span class="font-medium">Current schedule:</span>
                <br />
                <span class="font-medium">{{ schedule?.summary }}</span>
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
                v-if="isDarkMode"
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
              <Label for="dark-mode" class="text-base">Dark Mode</Label>
              <p class="text-sm text-muted-foreground">
                Switch between light and dark themes
              </p>
            </div>
            <Switch id="dark-mode" @click="toggleTheme" v-model="isDarkMode" />
          </div>
        </Card>

        <!-- Save Settings -->
        <Button class="w-full" @click="saveSettings"> Save Settings </Button>

        <!-- App Info Section -->
        <Card class="p-6 bg-transparent">
          <h3 class="font-medium mb-3">About Alfred</h3>
          <div class="space-y-2 text-sm text-muted-foreground">
            <p>Version {{ config.appVersion }}</p>
            <p>Built for family organization and management</p>
            <p>Telegram Mini App</p>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  ArrowLeft,
  Plus,
  Trash2,
  Eye,
  EyeOff,
  Copy,
  Check,
  LoaderCircle,
  HeartCrack,
  Frown,
} from "lucide-vue-next";

interface Note {
  id: number;
  key: string;
  value: string;
  isVisible: boolean;
  copyIcon: Component;
}

interface Secret {
  id: number | undefined;
  key: string;
  value: string;
  chatId: number;
  keyVersion: number;
  ivB64: string;
}

const route = useRoute();
const { public: config } = useRuntimeConfig();
const chatId: string = route.params.id as string;
const initDataRaw = useState<string>("initDataRaw");
const encryptionKey = useState<string>("encryptionKey");
const dek = await getDEK(encryptionKey.value);

const isDialogOpen = ref(false);
const isToastOpen = ref(false);
const toastStatus = ref("success");
const toastMessage = ref("");

const {
  data: notes,
  pending,
  error,
} = await useAsyncData<Note[]>(async () => {
  const res = await useFetch<Secret[]>(`/api/secrets/${chatId}`, {
    method: "GET",
    headers: {
      Authorization: `tma ${initDataRaw.value}`,
    },
  });
  console.log(res);

  let id = 1;
  let notes: Note[] = [];
  res.data.value?.forEach(async (s) => {
    const decrypted = await decrypt(dek, s.ivB64, s.value, chatId);
    notes.push({
      id: s.id ?? id++,
      key: s.key,
      value: decrypted,
      isVisible: false,
      copyIcon: Copy,
    });
  });

  return notes;
});

function back() {
  return navigateTo("/telegram");
}

function formatValue(value: string, isVisible: boolean) {
  return isVisible ? value : "*".repeat(15);
}

function toggleValueVisibility(id: number) {
  const note = notes.value?.find((n) => n.id === id);
  if (note) note.isVisible = !note.isVisible;
}

function setIsDialogOpen() {
  isDialogOpen.value = !isDialogOpen.value;
}

async function copyValue(id: number) {
  try {
    const note = notes.value?.find((n) => n.id === id);
    if (!note) return;

    const text = note.value;
    await navigator.clipboard.writeText(text);
    note.copyIcon = Check;

    isToastOpen.value = true;
    toastStatus.value = "success";
    toastMessage.value = "Copied!";

    setTimeout(() => {
      note.copyIcon = Copy;
    }, 1500); // Reset feedback after 1.5 seconds
  } catch (err) {
    console.error("Failed to copy text: ", err);
    toastStatus.value = "error";
    toastMessage.value = "Error copying to clipboard";
  }
}

const newKey = ref("");
const newValue = ref("");
async function submitNewNote() {
  const { iv, ciphertext } = await encrypt(
    dek,
    newValue.value,
    JSON.stringify(chatId)
  );

  const newSecret: Secret = {
    id: undefined,
    key: newKey.value,
    value: ciphertext,
    chatId: parseInt(chatId),
    keyVersion: config.keyVersion as number,
    ivB64: iv,
  };
  await $fetch("/api/secrets", {
    method: "POST",
    body: JSON.stringify(newSecret),
    headers: {
      Authorization: `tma ${initDataRaw.value}`,
    },
  }).catch((err) => {
    isDialogOpen.value = false;
    isToastOpen.value = true;
    toastStatus.value = "error";
    toastMessage.value = "An error occurred, failed to add new note";
    return;
  });

  isDialogOpen.value = false;
  isToastOpen.value = true;
  toastStatus.value = "success";
  toastMessage.value = "New note added!";
}
</script>

<template>
  <NuxtLink to="/telegram"></NuxtLink>
  <div class="min-h-screen bg-background">
    <!-- Header -->
    <div class="sticky top-0 bg-background border-b border-border p-4">
      <div class="flex items-center justify-between max-w-md mx-auto">
        <div class="flex items-center space-x-3">
          <Button @click="back" variant="ghost" size="sm" class="p-2">
            <ArrowLeft class="w-5 h-5" />
          </Button>
          <h1 class="text-xl font-medium">Notebook</h1>
        </div>

        <Dialog :open="isDialogOpen">
          <DialogTrigger as-child>
            <Button
              @click="setIsDialogOpen"
              size="sm"
              class="flex items-center space-x-2"
            >
              <Plus class="w-4 h-4" />
              <span>Add</span>
            </Button>
          </DialogTrigger>
          <DialogContent @dialog-close="setIsDialogOpen">
            <form @submit.prevent="submitNewNote">
              <DialogHeader>
                <DialogTitle>Add New Note</DialogTitle>
              </DialogHeader>
              <div class="space-y-4 pt-4">
                <div>
                  <Label htmlFor="key">Key</Label>
                  <Input
                    id="key"
                    placeholder="e.g., WiFi Password"
                    v-model="newKey"
                  />
                </div>
                <div>
                  <Label htmlFor="value">Value</Label>
                  <Input
                    id="value"
                    placeholder="e.g., MyPassword123"
                    v-model="newValue"
                  />
                </div>
                <div class="flex space-x-2 pt-2">
                  <Button class="flex-1" type="submit"> Add Note </Button>
                  <Button
                    variant="outline"
                    class="flex-1"
                    @click="setIsDialogOpen"
                    type="reset"
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </div>

    <!-- Notes List -->
    <div class="p-4">
      <div class="max-w-md mx-auto space-y-3">
        <Card v-for="note in notes" :key="note.id" class="p-4">
          <div class="flex items-start justify-between mb-2">
            <h3 class="font-medium text-foreground overflow-hidden">
              {{ note.key }}
            </h3>
            <Button
              variant="ghost"
              size="sm"
              class="text-destructive hover:text-destructive p-1"
            >
              <Trash2 class="w-4 h-4" />
            </Button>
          </div>
          <div class="flex items-center space-x-2 mb-2">
            <code
              class="flex-1 text-sm overflow-auto bg-muted p-2 rounded font-mono"
            >
              {{ formatValue(note.value, note.isVisible) }}
            </code>
            <Button
              @click="toggleValueVisibility(note.id)"
              variant="ghost"
              size="sm"
              class="p-2"
            >
              <EyeOff v-if="note.isVisible" class="w-4 h-4" />
              <Eye v-else class="w-4 h-4" />
            </Button>
            <Button
              @click="copyValue(note.id)"
              variant="ghost"
              size="sm"
              class="p-2"
            >
              <component :is="note.copyIcon" class="w-4 h-4" />
            </Button>
          </div>
        </Card>

        <!-- Error state -->
        <div v-if="error" class="flex flex-col items-center py-12">
          <Frown class="w-8 h-8 text-muted-foreground mb-2" />
          <p class="text-muted-foreground">An error occurred</p>
        </div>

        <!-- Pending state -->
        <div v-else-if="pending" class="flex flex-col items-center py-12">
          <LoaderCircle
            class="w-8 h-8 text-muted-foreground animate-spin mb-2"
          />
          <p class="text-muted-foreground">Loading</p>
        </div>

        <!-- Empty state -->
        <div v-else-if="!notes" class="text-center py-12">
          <p class="text-muted-foreground mb-4">No notes yet</p>
          <Button @click="setIsDialogOpen">
            <Plus class="w-4 h-4 mr-2" />
            Add your first note
          </Button>
        </div>
      </div>
    </div>

    <ToastRoot v-model:open="isToastOpen">
      <ToastDescription :status="toastStatus">{{
        toastMessage
      }}</ToastDescription>
    </ToastRoot>
  </div>
</template>

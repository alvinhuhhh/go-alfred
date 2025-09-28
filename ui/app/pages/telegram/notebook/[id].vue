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
  Frown,
} from "lucide-vue-next";

interface Note {
  id: number;
  key: string;
  value: string;
  isVisible: boolean;
  isDeleteDialogOpen: boolean;
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
const chatId: string = route.params.id as string;
const initDataRaw = useState<string>("initDataRaw");
const encryptionKey = useState<string>("encryptionKey");
const dek = await getDEK(encryptionKey.value);

const isDialogOpen = ref(false);
const isToastOpen = ref(false);
const toastStatus = ref("success");
const toastMessage = ref("");
const notes = ref<Note[]>([]);

const {
  data: raw,
  pending,
  error,
  refresh,
} = await useAsyncData<string>("notes", () =>
  $fetch(`/api/secrets/${chatId}`, {
    method: "GET",
    headers: {
      Authorization: `tma ${initDataRaw.value}`,
    },
  })
);

watch(
  raw,
  async (raw) => {
    if (!raw) return;
    const json = JSON.parse(raw);

    let id = 1;
    const arr = await Promise.all(
      json.map(async (s: Secret) => {
        const decrypted = await decrypt(dek, s.ivB64, s.value, chatId);
        return {
          id: s.id ?? id++,
          key: s.key,
          value: decrypted,
          isVisible: false,
          isDeleteDialogOpen: false,
          copyIcon: Copy,
        };
      })
    );

    notes.value = arr;
  },
  { immediate: true }
);

function back() {
  clearDialog();
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
  clearDialog();
  isDialogOpen.value = !isDialogOpen.value;
}

function setIsDeleteDialogOpen(noteId: number) {
  const note = notes.value?.find((n) => n.id === noteId);
  if (note) note.isDeleteDialogOpen = !note.isDeleteDialogOpen;
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
function clearDialog() {
  newKey.value = "";
  newValue.value = "";
}
async function submitNewNote() {
  const { iv, ciphertext } = await encrypt(dek, newValue.value, chatId);

  const newSecret: Secret = {
    id: undefined,
    key: newKey.value,
    value: ciphertext,
    chatId: parseInt(chatId),
    keyVersion: parseInt(import.meta.env.VITE_MASTER_KEY_VERSION) as number,
    ivB64: iv,
  };

  try {
    await $fetch("/api/secrets", {
      method: "POST",
      body: JSON.stringify(newSecret),
      headers: {
        Authorization: `tma ${initDataRaw.value}`,
      },
    });
  } catch (err) {
    clearDialog();
    isDialogOpen.value = false;
    isToastOpen.value = true;
    toastStatus.value = "error";
    toastMessage.value = "Failed to add new note";
    return;
  }

  clearDialog();
  isDialogOpen.value = false;
  isToastOpen.value = true;
  toastStatus.value = "success";
  toastMessage.value = "New note added!";
  await refresh();
}

async function deleteNote(noteId: number) {
  try {
    await $fetch(`/api/secrets/${noteId}`, {
      method: "DELETE",
      headers: {
        Authorization: `tma ${initDataRaw.value}`,
      },
    });
  } catch (err) {
    setIsDeleteDialogOpen(noteId);
    isToastOpen.value = true;
    toastStatus.value = "error";
    toastMessage.value = "Failed to delete note";
    return;
  }

  setIsDeleteDialogOpen(noteId);
  isToastOpen.value = true;
  toastStatus.value = "success";
  toastMessage.value = "Note deleted";
  await refresh();
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
        <div v-else-if="notes && notes.length === 0" class="text-center py-12">
          <p class="text-muted-foreground mb-4">No notes yet</p>
          <Button @click="setIsDialogOpen">
            <Plus class="w-4 h-4 mr-2" />
            Add your first note
          </Button>
        </div>

        <!-- Populate list -->
        <Card v-for="note in notes" :key="note.id" class="p-4">
          <div class="flex items-start justify-between mb-2">
            <h3 class="font-medium text-foreground overflow-hidden">
              {{ note.key }}
            </h3>

            <!-- Delete Dialog -->
            <Dialog :open="note.isDeleteDialogOpen">
              <DialogTrigger as-child>
                <Button
                  @click="setIsDeleteDialogOpen(note.id)"
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive p-1"
                >
                  <Trash2 class="w-4 h-4" />
                </Button>
              </DialogTrigger>
              <DialogContent @dialog-close="setIsDeleteDialogOpen(note.id)">
                <form @submit.prevent="deleteNote(note.id)">
                  <DialogHeader>
                    <DialogTitle>Delete Note</DialogTitle>
                  </DialogHeader>
                  <div class="space-y-4 py-4">
                    <p class="text-muted-foreground">
                      Are you sure you want to delete "{{ note.key }}"? This
                      action cannot be undone.
                    </p>
                  </div>
                  <div class="flex space-x-2 pt-2">
                    <Button variant="destructive" class="flex-1" type="submit"
                      >Delete</Button
                    >
                    <Button
                      variant="outline"
                      class="flex-1"
                      @click="setIsDeleteDialogOpen(note.id)"
                      type="reset"
                      >Cancel</Button
                    >
                  </div>
                </form>
              </DialogContent>
            </Dialog>
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
      </div>
    </div>

    <ToastRoot v-model:open="isToastOpen">
      <ToastDescription :status="toastStatus">{{
        toastMessage
      }}</ToastDescription>
    </ToastRoot>
  </div>
</template>

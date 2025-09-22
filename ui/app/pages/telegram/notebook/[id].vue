<script setup>
import {
  ArrowLeft,
  Plus,
  Trash2,
  Eye,
  EyeOff,
  Copy,
  Check,
} from "lucide-vue-next";

const isDialogOpen = ref(true);

const notes = ref([
  {
    id: 1,
    key: "CDC Vouchers",
    value: "www.link.com",
    isVisible: false,
    copyIcon: Copy,
  },
  {
    id: 2,
    key: "Wifi Password",
    value: "password@password",
    isVisible: false,
    copyIcon: Copy,
  },
]);

function back() {
  return navigateTo("/telegram");
}

function formatValue(value, isVisible) {
  return isVisible ? value : "*".repeat(15);
}

function toggleValueVisibility(id) {
  const note = notes.value.find((n) => n.id === id);
  note.isVisible = !note.isVisible;
}

function setIsDialogOpen() {
  isDialogOpen.value = !isDialogOpen.value;
}

async function copyValue(id) {
  try {
    const note = notes.value.find((n) => n.id === id);
    const text = note.value;
    await navigator.clipboard.writeText(text);
    note.copyIcon = Check;
    setTimeout(() => {
      note.copyIcon = Copy;
    }, 1500); // Reset feedback after 1.5 seconds
  } catch (err) {
    console.error("Failed to copy text: ", err);
    // Handle error, e.g., show an error message to the user
  }
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
          <DialogTrigger asChild>
            <Button
              @click="setIsDialogOpen"
              size="sm"
              class="flex items-center space-x-2"
            >
              <Plus class="w-4 h-4" />
              <span>Add</span>
            </Button>
          </DialogTrigger>
          <DialogContent
            @dialog-close="setIsDialogOpen"
            class="w-[90vw] max-w-md"
          >
            <DialogHeader>
              <DialogTitle>Add New Note</DialogTitle>
            </DialogHeader>
            <div class="space-y-4 pt-4">
              <div>
                <Label htmlFor="key">Key</Label>
                <Input id="key" placeholder="e.g., WiFi Password" />
              </div>
              <div>
                <Label htmlFor="value">Value</Label>
                <Input id="value" placeholder="e.g., MyPassword123" />
              </div>
              <div class="flex space-x-2 pt-2">
                <Button class="flex-1"> Add Note </Button>
                <Button
                  variant="outline"
                  @click="setIsDialogOpen"
                  class="flex-1"
                >
                  Cancel
                </Button>
              </div>
            </div>
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
              <EyeOff v-if="!note.isVisible" class="w-4 h-4" />
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

        <!-- Empty state -->
        <div v-if="notes.length < 1" class="text-center py-12">
          <p class="text-muted-foreground mb-4">No notes yet</p>
          <Button>
            <Plus class="w-4 h-4 mr-2" />
            Add your first note
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

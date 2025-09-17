<script setup>
import { init, retrieveLaunchParams } from "@telegram-apps/sdk-vue";
init();

const chatId = ref("unknown");

function setChatId() {
  const { tgWebAppData } = retrieveLaunchParams();
  console.log(tgWebAppData);

  if (!tgWebAppData) {
    return;
  }

  if (tgWebAppData.chat) {
    chatId.value = tgWebAppData.chat.id;
  } else if (tgWebAppData.user) {
    chatId.value = tgWebAppData.user.id;
  }
}

onMounted(() => {
  setChatId();
});
</script>

<template>
  <div
    class="h-screen w-screen flex flex-col bg-slate-950 justify-center items-center"
  >
    <h1 class="text-3xl text-slate-50">Welcome to Alfred!</h1>
    <body class="text-md text-slate-50">
      Chat ID: {{ chatId }}
    </body>
  </div>
</template>

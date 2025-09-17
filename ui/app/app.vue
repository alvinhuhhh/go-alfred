<script setup>
import { init, retrieveLaunchParams } from "@telegram-apps/sdk-vue";
init();
const { initDataRaw, initData, tgWebAppData } = retrieveLaunchParams();
console.log(initDataRaw);
console.log(initData);
console.log(tgWebAppData);

const chatId = ref("unknown");

if (tgWebAppData) {
  chatId.value = tgWebAppData.chat.id;
}
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

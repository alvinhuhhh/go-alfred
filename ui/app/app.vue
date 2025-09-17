<script setup>
import { retrieveLaunchParams } from "@telegram-apps/sdk";

const { initDataRaw, initData } = retrieveLaunchParams();
const chatId = ref("unknown");

console.log(initDataRaw, initData);
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

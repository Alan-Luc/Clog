<script setup lang="ts">
import Toaster from "./components/ui/toast/Toaster.vue";
import Login from "./components/views/Login.vue";
</script>

<template>
  <!-- <main> -->
  <RouterView />
  <Toaster />
  <!-- </main> -->
  <nav v-if="!$route.meta.hideNav">
    <RouterLink to="/login">Login</RouterLink>
    <RouterLink to="/about">about</RouterLink>
  </nav>
</template>

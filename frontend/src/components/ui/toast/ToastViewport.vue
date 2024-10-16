<script setup lang="ts">
import { cn } from "@/lib/utils";
import { ToastViewport, type ToastViewportProps } from "radix-vue";
import { computed, type HTMLAttributes } from "vue";

const props = defineProps<
  ToastViewportProps & { class?: HTMLAttributes["class"] }
>();

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props;

  return delegated;
});
</script>

<template>
  <ToastViewport
    v-bind="delegatedProps"
    :class="
      cn(
        'fixed top-0 left-[50%] z-[100] flex max-h-screen w-full translate-x-[-50%] flex-col-reverse p-4 sm:right-0 sm:flex-col md:max-w-[420px]',
        props.class,
      )
    "
  />
</template>

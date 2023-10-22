<style>
@tailwind base;
@tailwind components;
@tailwind utilities;

html {
  @apply font-sans;
}

html,
body,
#app {
  @apply w-full;
  height: 100dvh;
}

#app {
  @apply fixed overflow-auto flex flex-col justify-stretch items-stretch bg-base-200;
  height: 100dvh;
}

.btn,
.input,
.textarea,
.file-input,
.select,
.collapse,
.modal-box {
  @apply focus:outline-none;
}
</style>

<template>
  <div></div>
</template>

<script>
// import axios from "axios";

export default {
  data() {
    return {

    };
  },
  methods: {

  },
  async mounted() {

  },
};
</script>

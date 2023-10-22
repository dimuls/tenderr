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
  <div>
    <iframe class="border-0" src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=1" width="450" height="200"></iframe>
    <iframe class="border-0" src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=2" width="450" height="200"></iframe>
  </div>
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

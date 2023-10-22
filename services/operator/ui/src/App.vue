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
  min-height: 100dvh;
}

#app {
  @apply fixed overflow-auto flex flex-col justify-stretch items-stretch bg-base-200;
  min-height: 100dvh;
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
  <div class="flex w-full">
    <div class="flex flex-col">
      <div class="flex flex-col" v-for="ue in userErrors" :key="ue.id">
        <div>{{ elements[ue.elementId] }}</div>
        <div>{{ ue.message }}</div>
        <div>{{ ue.createdAt }}</div>
      </div>
    </div>
    <div class="flex flex-col">
      <iframe
        class="border-0"
        src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=1"
        width="450"
        height="200"
      ></iframe>
      <iframe
        class="border-0"
        src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=2"
        width="450"
        height="200"
      ></iframe>
    </div>
    <div class="flex flex-col"></div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      elements: {},
      userErrors: [],
      errorNotifications: [],
    };
  },
  methods: {},
  async mounted() {
    const [esRes, ueRes, enRes] = await Promise.all([
      axios.get("/elements"),
      axios.get("/user-errors"),
      axios.get("error-notifications"),
    ]);

    if (esRes.data) {
      for (const e of esRes.data) {
        this.elements[e.id] = e.name;
      }
    }

    if (ueRes.data) {
      this.userErrors.push(...ueRes.data);
    }

    if (enRes.data) {
      this.errorNotifications.push(...enRes.data);
    }
  },
};
</script>

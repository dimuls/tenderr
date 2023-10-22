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
  @apply flex flex-col justify-stretch items-stretch bg-base-200;
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
  <div class="flex w-full p-4 gap-8">
    <div class="flex-1 max-w-md flex flex-col">
      <div class="text-2xl p-4">Сообщенния об ошибках от пользователей</div>
      <div
        class="flex flex-col p-4 gap-2 bg-base-100"
        v-for="ue in userErrors"
        :key="ue.id"
      >
        <div>{{ elementsById[ue.elementId] }}</div>
        <div class="text-xl">{{ ue.message }}</div>
        <div class="text-right">{{ ue.createdAt }}</div>
      </div>
    </div>
    <div class="flex-1 flex flex-col gap-4">
      <div class="text-2xl px-4 pt-4">Метрики ошибок</div>
      <iframe
        class="border-0 flex-1 max-h-80"
        src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=1"
      ></iframe>
      <iframe
        class="border-0 flex-1 max-h-80"
        src="http://localhost:3000/d-solo/b4b3e695-8021-4b26-a5c1-8acef177e901/default?orgId=1&refresh=5s&panelId=2"
      ></iframe>
    </div>
    <div class="flex-1 max-w-md flex flex-col gap-4">
      <div class="flex flex-col gap-4">
        <div class="text-2xl px-4 pt-4">Новое оповещение об ошибке</div>
        <div class="flex flex-col gap-4 px-4 pb-4">
          <div class="form-control w-full max-w-xs">
            <label class="label">
              <span class="label-text">Элемент</span>
            </label>
            <select
              class="select select-bordered"
              v-model="newErrorNotification.elementId"
            >
              <option v-for="(e, i) in elements" :key="e.id" :value="e.id">
                {{ e.name }}
              </option>
            </select>
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Сообщение</span>
            </label>
            <textarea
              v-model="newErrorNotification.message"
              class="textarea textarea-bordered"
              rows="3"
            ></textarea>
          </div>
          <div class="flex justify-end">
            <button class="btn btn-ghost" @click="clearNewErrorNotification">
              Очистить
            </button>
            <button
              class="btn btn-primary"
              @click="publishNewErrorNotification"
            >
              Опубликовать
            </button>
          </div>
        </div>
      </div>
      <div class="text-2xl p-4">Оповещения об ошибках</div>
      <div
        class="flex flex-col p-4 gap-2 bg-base-100"
        v-for="en in errorNotifications"
        :key="en.id"
      >
        <div>{{ elementsById[en.elementId] }}</div>
        <div class="text-xl">{{ en.message }}</div>
        <div class="text-right">Опубликовано: {{ en.createdAt }}</div>
        <template v-if="!en.resolved">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Сообщение о решении проблемы</span>
            </label>
            <textarea
              v-model="en.resolveMessage"
              class="textarea textarea-bordered"
              rows="3"
            ></textarea>
          </div>
          <div class="flex justify-end mt-4">
            <button
              class="btn btn-primary"
              @click="resolveErrorNotification(en)"
            >
              Проблема решена
            </button>
          </div>
        </template>
        <div v-else class="text-right">Решено: {{ en.resolvedAt }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import dayjs from "dayjs";

const formatTime = (s) => {
  return dayjs(s).format("YYYY-MM-DD HH:mm:ss");
};

export default {
  data() {
    return {
      elementsById: {},
      elements: [],
      userErrors: [],
      errorNotifications: [],
      newErrorNotification: {
        elementId: "",
        message: "",
      },
    };
  },
  methods: {
    async reloadElements() {
      const { data } = await axios.get("/elements");
      this.elements.splice(0, this.elements.length);
      if (data) {
        this.elements.push(...data);
        for (const e of data) {
          this.elementsById[e.id] = e.name;
        }
      }
    },
    async reloadUserErrors() {
      const { data } = await axios.get("/user-errors");
      this.userErrors.splice(0, this.userErrors.length);
      if (data) {
        this.userErrors.push(
          ...data.map((ue) => ({
            ...ue,
            createdAt: formatTime(ue.createdAt),
          })),
        );
      }
    },
    async reloadErrorNotifications() {
      const { data } = await axios.get("/error-notifications");
      this.errorNotifications.splice(0, this.errorNotifications.length);
      if (data) {
        this.errorNotifications.push(
          ...data.map((en) => ({
            ...en,
            resolveMessage: "",
            createdAt: formatTime(en.createdAt),
            resolvedAt: en.resolvedAt ? formatTime(en.resolvedAt) : "",
          })),
        );
      }
    },
    async resolveErrorNotification(en) {
      await axios.patch("/error-notifications", {
        id: en.id,
        message: en.resolveMessage,
      });
      await this.reloadErrorNotifications();
    },
    clearNewErrorNotification() {
      this.newErrorNotification.elementId = "";
      this.newErrorNotification.message = "";
    },
    async publishNewErrorNotification() {
      await axios.post("/error-notifications", {
        elementId: this.newErrorNotification.elementId,
        message: this.newErrorNotification.message,
      });
      this.clearNewErrorNotification();
      await this.reloadErrorNotifications();
    },
  },
  async mounted() {
    await Promise.all([
      this.reloadElements(),
      this.reloadUserErrors(),
      this.reloadErrorNotifications(),
    ]);
  },
};
</script>

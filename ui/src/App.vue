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
  <div class="p-8 flex flex-col self-center w-full max-w-4xl">
    <div class="flex flex-col gap-4 w-full max-w-5xl">
      <div class="flex flex-col">
        <div class="flex gap-4">
          <div class="form-control w-48">
            <label class="label">
              <span class="label-text">Название</span>
            </label>
            <input
              type="text"
              class="input input-bordered"
              v-model="newClass.name"
            />
          </div>
          <div class="form-control w-full">
            <label class="label">
              <span class="label-text">Правила</span>
            </label>
            <textarea
              class="textarea textarea-bordered"
              rows="6"
              v-model="newClass.rules"
            ></textarea>
            <label class="label">
              <span class="label-text-alt"
                >Вводите каждое правило с новой строки</span
              >
            </label>
          </div>
        </div>
        <div class="w-full flex justify-end">
          <button class="btn btn-neutral" @click="addNewClass">Добавить</button>
        </div>
      </div>

      <div class="flex flex-col" v-for="c in classes" :key="c.id">
        <div class="flex gap-4">
          <div class="form-control w-48">
            <label class="label">
              <span class="label-text">Название</span>
            </label>
            <input type="text" class="input input-bordered" v-model="c.name" />
          </div>
          <div class="form-control w-full">
            <label class="label">
              <span class="label-text">Правила</span>
            </label>
            <textarea
              class="textarea textarea-bordered"
              rows="6"
              v-model="c.rules"
            ></textarea>
            <label class="label">
              <span class="label-text-alt"
                >Вводите каждое правило с новой строки</span
              >
            </label>
          </div>
        </div>
        <div class="w-full flex justify-end">
          <button class="btn btn-neutral" @click="removeClass(c)">
            Удалить
          </button>
        </div>
      </div>
      <div class="sticky bottom-0 pb-8 bg-base-200">
        <div class="divider" />
        <div class="flex gap-4 justify-end">
          <button class="btn btn-neutral" @click="clear">Сбросить</button>
          <button class="btn btn-neutral" @click="save">Сохранить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { TrashIcon } from "@heroicons/vue/24/outline/index.js";
import axios from "axios";

export default {
  components: {
    TrashIcon,
  },
  data() {
    return {
      origClasses: [],
      classes: [],
      newClass: {
        name: "",
        rules: "",
      },
      selectedClass: null,
    };
  },
  methods: {
    addNewClass() {
      if (!this.newClass.name || !this.newClass.rules) {
        return;
      }
      this.classes.unshift({
        name: this.newClass.name,
        rules: this.newClass.rules,
      });
      this.newClass.name = "";
      this.newClass.rules = "";
    },
    removeClass(c) {
      this.classes.splice(this.classes.indexOf(c), 1);
    },
    clear() {
      this.classes.splice(0, this.classes.length);
      this.classes.push(...this.origClasses);
    },
    async save() {
      const classes = [
        ...this.classes.map((c) => ({
          ...c,
          rules: c.rules.split(/\r?\n/),
        })),
      ];
      await axios.put("/classes", classes);
      this.origClasses.splice(0, this.origClasses.length);
      this.origClasses.push(...this.classes);
    },
  },
  async mounted() {
    const { data } = await axios.get("/classes");
    if (data) {
      this.origClasses.push(
        ...data.map((c) => ({
          ...c,
          rules: c.rules.join("\n"),
        })),
      );
      this.classes.push(...this.origClasses);
    }
  },
};
</script>

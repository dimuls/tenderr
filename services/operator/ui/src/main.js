import { createApp } from 'vue'
import App from './App.vue'
import axios from "axios";

axios.defaults.baseURL = import.meta.env.VITE_BASE_URL + "/api";
axios.defaults.timeout = 60000;

createApp(App).mount('#app')

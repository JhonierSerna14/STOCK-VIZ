import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/main.css"; // Asegúrate de que este archivo contenga las directivas de Tailwind

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount("#app");

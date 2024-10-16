import { createApp } from "vue";
import "./assets/index.css";
import "./assets/global.css";
import App from "@/App.vue";
import router from "./router/index.ts";

const app = createApp(App);

app.use(router);
app.mount("#app");

import { mount } from "svelte";
import "./app.css";
import App from "./App.svelte";

import "./wasm/wasm_exec";
import "./wasm/load";
import "./zxing";

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;

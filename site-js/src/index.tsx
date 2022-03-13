import { createSignal } from "solid-js";
import { render } from "solid-js/web";

function App() {
  const [message, setMessage] = createSignal<string>();

  (async () => {
    const resp = await fetch("/", { method: "POST" });
    const text = await resp.text();
    setMessage(text);
  })();

  return (
    <div>
      <h1>A cool webserver</h1>
      <div>
        This server is homegrown and runs on a kindle fire. It's written in 100%
        Go so it should be pretty fast.
      </div>
      <div>{message() !== undefined && message()}</div>
    </div>
  );
}

render(App, document.getElementById("root")!);

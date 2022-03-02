import { render } from "solid-js/web";

function App() {
  return (
    <div>
      <h1>A cool webserver</h1>
      <div>
        This server is homegrown and runs on a kindle fire. It's written in 100%
        Go so it should be pretty fast.
      </div>
    </div>
  );
}

render(App, document.getElementById("root")!);

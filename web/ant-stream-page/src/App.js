import { useState } from "react";
import WebRTCConfigure from "./Containers/WebRTCConfigure";
import WSConfigure from "./Containers/WSConfigure";

function App() {
  return (
      <div>
        <WSConfigure />
        <WebRTCConfigure />
      </div>
  );
}

export default App;

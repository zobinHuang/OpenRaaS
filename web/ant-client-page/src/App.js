import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from './Containers/HomePage';
import InfoConfig from "./Configurations/InfoConfig.json"
import APIConfig from "./Configurations/APIConfig.json"
import SigninPage from "./Containers/SigninPage";
import SignupPage from "./Containers/SignupPage";
import { useState } from "react";
import SnackBar from "./Components/SnackBar";
import UserPage from "./Containers/UserPage";
import Backdrop from '@mui/material/Backdrop';
import CircularProgress from '@mui/material/CircularProgress';
import { useSelector, useDispatch } from 'react-redux'
import WebsocketCallback from "./Components/WebsocketCallback/ws_callback";
import VideoStreamPage from "./Containers/VideoStreamPage";

/*
    component: App
    description: antcloud web application
*/
function App() {
  // get global state of backdrop
  const backdropEnabled = useSelector(state => state.backdrop.backdropEnabled)

  /*
    @state: terminalWsMap
    @description:
        mapping relationship of "terminal frontent id <-> websocket"
  */
  const [terminalWsMap, setTerminalWsMap] = useState(new Map())

  /*
    @state: terminalRtcPeerMap
    @description:
        mapping relationship of "terminal frontent id <-> RtcPeer"
  */
  const [terminalRtcPeerMap, setTerminalRtcPeerMap] = useState(new Map())

  /*
    @state: terminalDynamicState
    @description:
        mapping relationship of "terminal frontent id <-> terminalDynamicState"
  */
  const [terminalDynamicState, setTerminalDynamicState] = useState(new Map())



  return (
    <BrowserRouter>
      {/* Websocket Callback Registeration */}
      <WebsocketCallback 
        terminalWsMap={terminalWsMap}
        setTerminalWsMap={setTerminalWsMap}
        terminalRtcPeerMap={terminalRtcPeerMap}
        setTerminalRtcPeerMap={setTerminalRtcPeerMap}
        terminalDynamicState={terminalDynamicState}
        setTerminalDynamicState={setTerminalDynamicState}
      />

      {/* Page Router */}
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/signin" element={<SigninPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/stream" element={<VideoStreamPage 
          terminalRtcPeerMap={terminalRtcPeerMap}
          setTerminalRtcPeerMap={setTerminalRtcPeerMap}
          terminalDynamicState={terminalDynamicState}
        />} />
      </Routes>

      {/* SnackBar Component */}
      <SnackBar />

      {/* Backdrop Component */}
      <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={backdropEnabled}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
    </BrowserRouter>
  );
}

export default App;
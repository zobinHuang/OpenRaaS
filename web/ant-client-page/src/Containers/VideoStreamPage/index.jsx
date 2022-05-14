import React, { useRef, useEffect } from 'react';
import styled from 'styled-components';
import { useDispatch, useSelector } from 'react-redux';
import { useSearchParams } from "react-router-dom"
import PubSub from 'pubsub-js';

const VideoStreamContainer = styled.div`
`

const VideoShowcaseContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const VideoStreamPage = (props) => {
    // get props
    const {terminalRtcPeerMap, setTerminalRtcPeerMap} = props;

    // create ref hook for stream
    const streamRef = useRef(null)
    
    // get global state
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    const CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta

    // obtain terminal key and peer connection object from url and global map respectively
    const [searchParams, setSearchParams] = useSearchParams();
    const terminalKey = searchParams.get("key")
    let RtcPeer = terminalRtcPeerMap.get(terminalKey)

    // add track
    useEffect(() => {
        streamRef.current.srcObject = RtcPeer.mediaStream[0]
        console.log(RtcPeer)
    },[])
    
    return <VideoStreamContainer>
        <VideoShowcaseContainer><video 
            style={{
                width: CurrentSelectedTerminal.screenWidth,
                height: CurrentSelectedTerminal.screenHeight,  
                margin: 5,
                backgroundColor: "#000000",
            }}
            autoPlay
            ref={streamRef}
        /></VideoShowcaseContainer>
    </VideoStreamContainer>
}

export default VideoStreamPage
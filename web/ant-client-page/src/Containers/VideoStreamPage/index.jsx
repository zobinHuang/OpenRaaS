import React, { useRef, useEffect } from 'react';
import styled from 'styled-components';
import { useDispatch, useSelector } from 'react-redux';
import { useSearchParams } from "react-router-dom"
import PubSub from 'pubsub-js';

const VideoStreamContainer = styled.div`
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const VideoShowcaseContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const EVENT_KEYDOWN = "KEYDOWN"
const EVENT_KEYUP = "KEYUP"
const EVENT_MOUSEDOWN = "MOUSEDOWN"
const EVENT_MOUSEUP = "MOUSEUP"
const EVENT_MOUSEMOVE = "MOUSEMOVE"

const MOUSE_LEFT = 0
const MOUSE_RIGHT = 1

const VideoStreamPage = (props) => {
    // get props
    const {terminalRtcPeerMap, setTerminalRtcPeerMap, terminalDynamicState} = props;

    // create ref hook for stream
    const streamRef = useRef(null)
    
    // obtain terminal key and peer connection object from url and global map respectively
    const [searchParams, setSearchParams] = useSearchParams();
    const terminalKey = searchParams.get("key")
    let RtcPeer = terminalRtcPeerMap.get(terminalKey)
    let instanceDynamicState = terminalDynamicState.get(terminalKey)

    // get global state
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const CurrentSelectedTerminal = StateTerminals.terminalsMap[terminalKey]
    
    // add track
    useEffect(() => {
        streamRef.current.srcObject = RtcPeer.mediaStream[0]
        
        // console.log('test')
        console.log(RtcPeer.PeerConnection)
        const xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://kb109.dynv6.net:52109/api/scheduler/record_history', true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        const latencies = [];

        const handleRTCStats = async () => {
            try {
                if (RtcPeer.PeerConnection) {
                    const stats = await RtcPeer.PeerConnection.getStats();

                    stats.forEach(report => {
                        // console.log(report.type,report)
                        if (report.type === "candidate-pair" && report.lastPacketReceivedTimestamp != undefined && report.lastPacketSentTimestamp != undefined) {
                            console.log(report.type,report)
                            const latency = report.totalRoundTripTime / report.responsesReceived * 1000
                            // const latency = report.lastPacketReceivedTimestamp - report.lastPacketSentTimestamp
                            // const jitter = report.jitter * 1000;

                            if (!isNaN(latency)) {
                                latencies.push(latency); 
                            }


                        }
                    });

                }
            } catch (error) {
                console.error("RTC Info error:", error);
            }
        };

        const sendAverageLatency = () => {
            if (latencies.length > 0) {
                console.log(latencies)
                // Remove minimum and maximum values
                const sortedLatencies = latencies.sort((a, b) => a - b);
                const trimmedLatencies = sortedLatencies.slice(2, -2);

                const averageLatency = trimmedLatencies.reduce((sum, latency) => sum + latency, 0) / latencies.length;

                const jsonData = {
                    instance_id: String(instanceDynamicState.instanceSchedulerID),
                    latency: String(averageLatency) + 'ms',
                };

                console.log(jsonData);
                xhr.send(JSON.stringify(jsonData));
                xhr.onreadystatechange = function() {
                    if (xhr.readyState === 4 && xhr.status === 200) {
                        const response = JSON.parse(xhr.responseText);
                        console.log(response);
                    }
                }

                latencies.length = 0;
            }
        };

        const startStatsInterval = () => {
            const statsInterval = setInterval(handleRTCStats, 1000);
          
            setTimeout(() => {
                clearInterval(statsInterval); // 停止每秒记录
                sendAverageLatency(); // 计算平均值并发送
            }, 20000);
        };

        // 延迟10秒后启动 statsInterval
        setTimeout(startStatsInterval, 0);

        //console.log('end')

        /*
            @callback: keydown
            @description: send keydown event and corresponding key code to remote peer
        */
        document.addEventListener("keydown", (event) => {
            console.log(event.keyCode)
            RtcPeer.inputChannel.send(JSON.stringify({
                packet_type: EVENT_KEYDOWN,
                data: JSON.stringify({
                  key_code: event.keyCode,
                }),
            }))
        })

        /*
            @callback: keydown
            @description: send keyup event and corresponding key code to remote peer
        */
        document.addEventListener("keyup", (event) => {
            console.log(event.keyCode)
            RtcPeer.inputChannel.send(JSON.stringify({
                packet_type: EVENT_KEYUP,
                data: JSON.stringify({
                    key_code: event.keyCode,
                }),
            }))
        })

        /*
            @callback: mousedown
            @description: send mousedown event and corresponding metadata to remote peer
        */
        document.addEventListener("mousedown", (event) => {
            console.log(event.button)

            let boundRect = streamRef.current.getBoundingClientRect()

            RtcPeer.inputChannel.send(JSON.stringify({
                packet_type: EVENT_MOUSEDOWN,
                data: JSON.stringify({
                    width: boundRect.width,
                    height: boundRect.height,
                    x: event.offsetX,
                    y: event.offsetY,
                    is_left: event.button === MOUSE_LEFT ? 1 : 0
                }),
            }))
        })

        /*
            @callback: contextmenu
            @description: diable right click button
        */
        document.addEventListener("contextmenu", (event) => {
            event.preventDefault()
        })

        /*
            @callback: mouseup
            @description: send mouseup event and corresponding metadata to remote peer
        */
        document.addEventListener("mouseup", (event) => {
            console.log(event.button)

            let boundRect = streamRef.current.getBoundingClientRect()

            RtcPeer.inputChannel.send(JSON.stringify({
                packet_type: EVENT_MOUSEUP,
                data: JSON.stringify({
                    width: boundRect.width,
                    height: boundRect.height,
                    x: event.offsetX,
                    y: event.offsetY,
                    is_left: event.button === MOUSE_LEFT ? 1 : 0
                }),
            }))
        })

        /*
            @callback: mousemove
            @description: send mousemove event and corresponding metadata to remote peer
        */
        document.addEventListener("mousemove", (event) => {
            // the log will interfere the user experience
            // console.log(event.button)

            let boundRect = streamRef.current.getBoundingClientRect()

            RtcPeer.inputChannel.send(JSON.stringify({
                packet_type: EVENT_MOUSEMOVE,
                data: JSON.stringify({
                    width: boundRect.width,
                    height: boundRect.height,
                    x: event.offsetX,
                    y: event.offsetY,
                    is_left: event.button === MOUSE_LEFT ? 1 : 0
                }),
            }))
        })

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
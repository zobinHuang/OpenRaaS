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
    const {terminalRtcPeerMap, setTerminalRtcPeerMap} = props;

    // create ref hook for stream
    const streamRef = useRef(null)
    
    // obtain terminal key and peer connection object from url and global map respectively
    const [searchParams, setSearchParams] = useSearchParams();
    const terminalKey = searchParams.get("key")
    let RtcPeer = terminalRtcPeerMap.get(terminalKey)

    // get global state
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const CurrentSelectedTerminal = StateTerminals.terminalsMap[terminalKey]

    // add track
    useEffect(() => {
        streamRef.current.srcObject = RtcPeer.mediaStream[0]
        
        console.log('test')
        console.log(RtcPeer.PeerConnection)
        //console.log(typeof RtcPeer);
        //console.log(RtcPeer instanceof RTCPeerConnection);

        
        // 创建一个处理RTC统计信息的函数
        const handleRTCStats = async () => {
            try {
                if (RtcPeer.PeerConnection) {
                    const stats = await RtcPeer.PeerConnection.getStats();
                    stats.forEach(report => {
                        //console.log(report)
                        if (report.type === "inbound-rtp" && report.kind === "video") {
                            // 获取RTCP的时间戳
                            const rtcpTimestamp = report.lastPacketReceivedTimestamp;
                            // 计算延迟（以毫秒为单位）
                            const now = new Date().getTime();
                            const latency = now - rtcpTimestamp;
                            //const latency = report.totalInterFrameDelay * 1000; // 转换为毫秒
                            const jitter = report.jitter * 1000; // 转换为毫秒
                            console.log("Latency:", latency, "ms");
                            console.log("Jitter:", jitter, "ms");
                            
                        }
                    });
                }
            } catch (error) {
                console.error("RTC Info error:", error);
            }
        };

        // 添加定时器以定期获取RTC统计信息
        const statsInterval = setInterval(handleRTCStats, 1000); // 每秒获取一次统计信息


        console.log('end')

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
            console.log(event.button)

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
        // 在组件卸载时清除定时器以防止内存泄漏
        return () => {
            clearInterval(statsInterval);
        };

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
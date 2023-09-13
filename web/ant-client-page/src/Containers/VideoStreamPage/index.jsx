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
        
        // ����һ������RTCͳ����Ϣ�ĺ���
        const handleRTCStats = async () => {
            try {
                const stats = await RtcPeer.peerConnection.getStats();
                stats.forEach(report => {
                    if (report.type === "inbound-rtp" && report.kind === "video") {
                        const latency = report.roundTripTime * 1000; // ת��Ϊ����
                        const jitter = report.jitter * 1000; // ת��Ϊ����
                        console.log("�ӳ٣�Latency��:", latency, "����");
                        console.log("������Jitter��:", jitter, "����");
                        
                    }
                });
            } catch (error) {
                console.error("��ȡRTCͳ����Ϣʱ����:", error);
            }
        };

        // ��Ӷ�ʱ���Զ��ڻ�ȡRTCͳ����Ϣ
        const statsInterval = setInterval(handleRTCStats, 1000); // ÿ���ȡһ��ͳ����Ϣ


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
        // �����ж��ʱ�����ʱ���Է�ֹ�ڴ�й©
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
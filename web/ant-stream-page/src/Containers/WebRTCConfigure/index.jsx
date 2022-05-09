import { useEffect, useState } from "react";
import { useSelector, useDispatch } from 'react-redux'
import PubSub from 'pubsub-js';

const WebRTCConfigure = (props) => {
    /*
        @callback: SetConsumerType
        @description:
            set consumer type as "stream"
    */
    const SetConsumerType = (msg, payload) => {
        const reqWSPacket = JSON.stringify({
            packet_type: "init_consumer_type",
            data: JSON.stringify({ 
                consumer_type: "stream"
            }),
        })
        payload.Socket.send(reqWSPacket);
    }

    const [webRTCConnection, setWebRTCConnection] = useState(null)

    /*
        @callback: InitWebRTC
        @description:
            initialize webrtc
    */
    const InitWebRTC = (msg, payload) => {
        // create peer connection
        let connection = new RTCPeerConnection({
            iceServers: JSON.parse(payload.WSPacket.data.iceservers),
        });
        setWebRTCConnection(connection)
    }
   
    useEffect(() => {
        var token_websocket_opened = PubSub.subscribe("websocket_opened", SetConsumerType)
        var token_webrtc_init_start = PubSub.subscribe("webrtc_init_start", InitWebRTC)
        return ()=>{
            PubSub.unsubscribe(token_websocket_opened)
            PubSub.unsubscribe(token_webrtc_init_start)
        }
    }, [])

    return null
}

export default WebRTCConfigure
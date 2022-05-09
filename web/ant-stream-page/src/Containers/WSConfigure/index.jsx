import { useEffect } from "react";
import { useSelector, useDispatch } from 'react-redux'
import PubSub from 'pubsub-js';

const WSConfigure = (props) => {
  const StateInfo = useSelector(state => state.info.StateInfo)
  const StateAPI = useSelector(state => state.api.StateAPI)

  /*
      @function: initilizeWebSocket
      @description:
          initialize websocket connection, and config websocket message router
  */
  const initilizeWebSocket = () => {
      // initilize websocket connection
      const ws = new WebSocket(`${StateAPI.WSProtocol}://${StateAPI.AuthHostAddr}:${StateAPI.WSPort}${StateAPI.WSBaseURL}/${StateAPI.WSAPI.WebSocketConnect}?type=${StateInfo.ClientType}`)

      // event: websocket opened
      ws.onopen = (event) => {
        console.log("Info - web socket opened", {event})
        PubSub.publish('websocket_opened', {Socket: ws});
      }
      
      // event: websocket closed
      ws.onclose = (event) => {
        console.log("Info - web socket closed", {event})
        PubSub.publish('websocket_closed', {});
      }
      
      // event: received websocket packet
      ws.onmessage = (event) => {
        const wsPacket = JSON.parse(event.data)
        
        // extract packet data
        const wsPacketData = JSON.parse(wsPacket.data);
        wsPacket.data = wsPacketData
        console.log("Info - web socket message: ", {wsPacket})

        // publish events
        switch(wsPacket.packet_type){
          /*
              @ case: webrtc_init_start
              @ description: notification of starting webrtc initialization process
           */
          case "webrtc_init_start":
            PubSub.publish('webrtc_init_start', { 
              Socket: ws, 
              WSPacket: wsPacket 
            });
            break
          
          default:
            console.log("Warn - receive unknown packet type: ", wsPacket.packet_type)
        }
      }
  
      ws.onerror = (error) => {
        console.log("Warn - web socket error: ", {error})
        PubSub.publish('websocket_error', {});
      }
  }

  /*
      @callback: webSocketKeepAlive
      @description:
          periodically send heartbeat to server (every 10 second)
  */
  const webSocketKeepAlive = (msg, payload) => {
    const reqWSPacket = JSON.stringify({
        packet_type: "keep_consumer_alive",
    })
    console.log("Info - start periodically heartbeat")
    setInterval(() => payload.Socket.send(reqWSPacket), 10000)
  }

  useEffect(() => {
      // initailize websocket
      initilizeWebSocket()

      // maintain heartbeat
      PubSub.subscribe('websocket_opened', webSocketKeepAlive)
  }, [])

  return null
}

export default WSConfigure;
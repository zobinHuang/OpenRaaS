import React, {useEffect, useState} from "react";
import PubSub from 'pubsub-js';
import GetTimestamp from '../../Utils/get_timestamp';
import { actions as TerminalActions } from '../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import { useDispatch, useSelector } from 'react-redux';

const WebsocketCallback = (props) => {
    const dispatch = useDispatch()

    // get global state
    const StateTerminals = useSelector(state => state.terminal)

    /*
        @state: "terminal <-> websocket" mapping relationship
        @description:
    */
    const [terminalWsMap, setTerminalWsMap] = useState(new Map())

    /*
        @callback: callback_registerTerminalWebsocket
        @description: 
            callback function for registering new "terminal <-> websocket" mapping relationship
    */
    const callback_registerTerminalWebsocket = (msg, payload) => {
        setTerminalWsMap(terminalWsMap.set(payload.TerminalKey, payload.Websocket));
    }

    /*
        @callback: callback_configWebsocketState
        @description: 
            callback function for registering websocket behavior under different state
    */
    const callback_configWebsocketState = (msg, payload) => {
        // fetch payload
        let ws = payload.Websocket
        let terminalKey = payload.TerminalKey
        let terminalName = payload.TerminalName
        let stateTerminals = payload.StateTerminals
        
        // record index of websocket keep alive interval
        let WSKeepAliveIndex = null
        
        /*
            @callback: onopen
            @description:
                invoked when websocket opened
        */
        ws.onopen = (event) => {
            // store websocket object
            PubSub.publish('register_terminal_websocket', {
                TerminalKey: terminalKey,
                Websocket: ws
            });

            // append terminal log
            dispatch(TerminalActions.updateTerminal({
                "type": "APPEND_LOG_CONTENT",
                "terminal_key": `${terminalKey}`,
                "log_priority": "SUCCESS",
                "log_time": GetTimestamp(),
                "log_content": "successfully connect to scheduler",
            }))

            // start periodically heartbeat
            const keepAliveWSPacket = JSON.stringify({
                packet_type: "keep_consumer_alive",
            })
            console.log("Info - start periodically heartbeat")
            WSKeepAliveIndex = setInterval(() => ws.send(keepAliveWSPacket), 10000)

            // append terminal log
            dispatch(TerminalActions.updateTerminal({
                "type": "APPEND_LOG_CONTENT",
                "terminal_key": `${terminalKey}`,
                "log_priority": "INFO",
                "log_time": GetTimestamp(),
                "log_content": "start websocket keep alive heartbeat",
            }))

            // publish event
            PubSub.publish('websocket_opened', {Socket: ws});

            // change current step
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_STEP",
                "terminal_key": `${terminalKey}`,
                "current_step_index": 2
            }))

            // send consummer metadata to scheduler
            let reqWSPacket = JSON.stringify({
                packet_type: "init_consumer_metadata",
                data: JSON.stringify({ 
                    consumer_type: stateTerminals.terminalsMap[terminalKey].applicationMeta.currentSelectedApplicationType,
                }),
            })
            ws.send(reqWSPacket)

            // append terminal log
            dispatch(TerminalActions.updateTerminal({
                "type": "APPEND_LOG_CONTENT",
                "terminal_key": `${terminalKey}`,
                "log_priority": "INFO",
                "log_time": GetTimestamp(),
                "log_content": "send consumer metadata to scheduler",
            }))
        }

        /*
            @callback: onmessage
            @description:
                invoked when recv message from websocket
        */
        ws.onmessage = (event) => {
            const wsPacket = JSON.parse(event.data)

            // extract packet data
            const wsPacketData = JSON.parse(wsPacket.data);
            wsPacket.data = wsPacketData
            console.log("Info - web socket message: ", {wsPacket})

            // publish events
            switch(wsPacket.packet_type){
                /*
                    @ case: notify_ice_server
                    @ description: notification the ice server for webrtc initilization
                */
                case "notify_ice_server":
                    PubSub.publish('notify_ice_server', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;
                
                /*
                    @ case: state_failed_provider_schedule
                    @ description: notification that scheduler failed to conduct schedule
                */
                case "state_failed_provider_schedule":
                    PubSub.publish('state_failed_provider_schedule', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;
                
                /*
                    @ case: state_provider_scheduled
                    @ description: notification that provider has selected
                */
                case "state_provider_scheduled":
                    PubSub.publish('state_provider_scheduled', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;

                /*
                    @ case: unknown websocket packet type
                    @ description: prompt unknown websocket packet type
                */
                default:
                    console.log("Warn - receive unknown packet type: ", wsPacket.packet_type)
            }
        }

        /*
            @callback: onclose
            @description:
                invoked when websocket closed
        */
        ws.onclose = (event) => {
            // log
            console.log("Info - web socket closed", {event})

            // show snack bar
            dispatch(SnackBarActions.showSnackBar(`Websocket of terminal for "${terminalName}" lose its connection`))

            // clear keep alive interval
            clearInterval(WSKeepAliveIndex);

            // publish event
            PubSub.publish('websocket_closed', {});

            // append terminal log
            dispatch(TerminalActions.updateTerminal({
                "type": "APPEND_LOG_CONTENT",
                "terminal_key": `${terminalKey}`,
                "log_priority": "WARN",
                "log_time": GetTimestamp(),
                "log_content": "websocket connection closed",
            }))

            // change current step
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_STEP",
                "terminal_key": `${terminalKey}`,
                "current_step_index": 0
            }))
        }

        /*
            @callback: onerror
            @description:
                invoked error occurs on the websocket
        */
        ws.onerror = (error) => {
            console.log("Warn - web socket error: ", {error})

            // append terminal log
            dispatch(TerminalActions.updateTerminal({
                "type": "APPEND_LOG_CONTENT",
                "terminal_key": `${terminalKey}`,
                "log_priority": "ERROR",
                "log_time": GetTimestamp(),
                "log_content": `${error}`,
            }))

            PubSub.publish('websocket_error', {});
        }
    }

    /*
        @callback: callback_deleteTerminalWebsocket
        @description: 
            callback function for closing websocket and deleting specified "terminal <-> websocket" 
            mapping relationship
    */
    const callback_deleteTerminalWebsocket = (msg, payload) => {
        if(terminalWsMap.get(payload.TerminalKey) !== undefined){
            terminalWsMap.get(payload.TerminalKey).close()
            setTerminalWsMap(delete terminalWsMap[payload.TerminalKey]);
        } else {
            console.log(`Warn - unknown terminal key '${payload.TerminalKey}' while deleting "terminal-websocket" mapping, abandoned`)
        }
    }

    /*
        @callback: callback_recvNotifyIceServer
        @description: callback function for receiving websocket packet of type "notify_ice_server"
    */
    const callback_recvNotifyIceServer = (msg, payload) => {
        // append terminal log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `obtain ice servers: ${payload.WSPacket.data.iceservers}`,
        }))

        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `obtain client id: ${payload.WSPacket.data.client_id}`,
        }))

        // store iceserver
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_ICESERVERS",
            "terminal_key": `${payload.TerminalKey}`,
            "ice_servers": `${payload.WSPacket.data.iceservers}`
        }))

        // store client id from scheduler
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_CLIENT_ID",
            "terminal_key": `${payload.TerminalKey}`,
            "client_id": `${payload.WSPacket.data.client_id}`
        }))

        // send metadata of selected application to scheduler
        let reqWSPacket = JSON.stringify({
            packet_type: "select_stream_application",
            data: JSON.stringify({ 
                application_id: payload.StateTerminals.terminalsMap[payload.TerminalKey].applicationMeta.currentSelectedApplication.id,
            }),
        })
        payload.Socket.send(reqWSPacket)
    }

    /*
        @callback: callback_stateFailedProviderSchedule
        @description: callback function for state notification of failed to find proper provider
    */
    const callback_stateFailedProviderSchedule = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "ERROR",
            "log_time": GetTimestamp(),
            "log_content": `failed to schedule provider: ${payload.WSPacket.data.error}`,
        }))

        // close websocket
        PubSub.publish('delete_terminal_websocket', {
            TerminalKey: payload.TerminalKey
        });

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": 0
        }))

        // unconfirm websocket connection started
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_WS_CONNECTION_STARTED",
            "terminal_key": `${payload.TerminalKey}`,
        }))
    }
    
    /*
        @callback: callback_stateProviderScheduled
        @description: callback function for state notification of found proper provider
    */
    const callback_stateProviderScheduled = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "SUCCESS",
            "log_time": GetTimestamp(),
            "log_content": `scheduler has found one provider to serve, provider id: ${payload.WSPacket.data.provider_id}`,
        }))

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": 3
        }))
    }

    /*
        @function: registerRecvCallback
        @description: register recv callback functions of given websocket
    */
    const registerRecvCallback = () => {
        PubSub.subscribe('register_terminal_websocket', callback_registerTerminalWebsocket)
        PubSub.subscribe('config_websocket_state', callback_configWebsocketState)
        PubSub.subscribe('delete_terminal_websocket', callback_deleteTerminalWebsocket)
        PubSub.subscribe('notify_ice_server', callback_recvNotifyIceServer)
        PubSub.subscribe('state_failed_provider_schedule', callback_stateFailedProviderSchedule)
        PubSub.subscribe('state_provider_scheduled', callback_stateProviderScheduled)
    }

    useEffect( () => {
        registerRecvCallback()
        console.log("register ws receive callback!")
    },[])

    return <div />
}

export default WebsocketCallback
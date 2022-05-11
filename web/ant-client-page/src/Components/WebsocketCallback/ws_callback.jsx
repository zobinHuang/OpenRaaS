import React, {useEffect, useState} from "react";
import PubSub from 'pubsub-js';
import GetTimestamp from '../../Utils/get_timestamp';
import { actions as TerminalActions } from '../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import { useDispatch, useSelector } from 'react-redux';
import { TERMINAL_STEP_SCHEDULE_COMPUTE_NODE, TERMINAL_STEP_CONFIG_INSTANCE, TERMINAL_STEP_SCHEDULE_STORAGE_NODE, TERMINAL_STEP_PREPARE_INSTANCE, TERMINAL_STEP_RUN_INSTANCE } from "../../Containers/UserPage/terminals";

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
                "current_step_index": TERMINAL_STEP_SCHEDULE_COMPUTE_NODE
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
                    @ case: state_failed_select_storage
                    @ description: notification that provider failed to find proper storage nodes
                */
                case "state_failed_select_storage":
                    PubSub.publish('state_failed_select_storage', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;
                
                /*
                    @ case: state_selected_storage
                    @ description: notification that storage nodes has selected
                */
                case "state_selected_storage":
                    PubSub.publish('state_selected_storage', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;
                
                /*
                    @ case: state_run_instance
                    @ description: notification that provider successfully runs instance
                */
                case "state_run_instance":
                    PubSub.publish('state_run_instance', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;

                /*
                    @ case: state_failed_run_instance
                    @ description: notification that provider failed to run instance
                */
                case "state_failed_run_instance":
                    PubSub.publish('state_failed_run_instance', { 
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
                "current_step_index": TERMINAL_STEP_CONFIG_INSTANCE,
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
                screen_height: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].screenHeight}`,
                screen_width: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].screenWidth}`,
                application_fps: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].currentFPS}`,
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
            "current_step_index": TERMINAL_STEP_CONFIG_INSTANCE
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
            "current_step_index": TERMINAL_STEP_SCHEDULE_STORAGE_NODE
        }))
    }

    /*
        @callback: callback_stateFailedStorageSchedule
        @description: callback function for state notification of failed to find proper storage nodes
    */
    const callback_stateFailedStorageSchedule = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "ERROR",
            "log_time": GetTimestamp(),
            "log_content": `provider failed to find proper storage node: ${payload.WSPacket.data.error}`,
        }))

        // close websocket
        PubSub.publish('delete_terminal_websocket', {
            TerminalKey: payload.TerminalKey
        });

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": TERMINAL_STEP_CONFIG_INSTANCE
        }))

        // unconfirm websocket connection started
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_WS_CONNECTION_STARTED",
            "terminal_key": `${payload.TerminalKey}`,
        }))
    }

    /*
        @callback: callback_stateStorageScheduled
        @description: callback function for state notification of found proper storage nodes
    */
    const callback_stateStorageScheduled = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "SUCCESS",
            "log_time": GetTimestamp(),
            "log_content": `provider has found proper storage nodes: depository address ${payload.WSPacket.data.target_depository}, filestore address ${payload.WSPacket.data.target_filestore}`,
        }))

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": TERMINAL_STEP_PREPARE_INSTANCE
        }))
    }

    /*
        @callback: callback_stateInstanceRunning
        @description: callback function for state notification of instance is now running
    */
    const callback_stateInstanceRunning = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "SUCCESS",
            "log_time": GetTimestamp(),
            "log_content": `provider is now successfully running instance with instance id ${payload.WSPacket.data.stream_instance_id}`,
        }))

        // update instance id in scheduler
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_INSTANCE_SCHEDULER_ID",
            "terminal_key": `${payload.TerminalKey}`,
            "instance_scheduler_id": `${payload.WSPacket.data.stream_instance_id}`
        }))

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": TERMINAL_STEP_RUN_INSTANCE
        }))
    }

    /*
        @callback: callback_stateFailedInstanceRunning
        @description: callback function for state notification of instance is failed to run
    */
    const callback_stateFailedInstanceRunning = (msg, payload) => {
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "ERROR",
            "log_time": GetTimestamp(),
            "log_content": `provider failed to run instance: ${payload.WSPacket.data.error}`,
        }))

        // close websocket
        PubSub.publish('delete_terminal_websocket', {
            TerminalKey: payload.TerminalKey
        });

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${payload.TerminalKey}`,
            "current_step_index": TERMINAL_STEP_CONFIG_INSTANCE
        }))

        // unconfirm websocket connection started
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_WS_CONNECTION_STARTED",
            "terminal_key": `${payload.TerminalKey}`,
        }))
    }

    /*
        @callback: callback_InitializeWebRTCConnection
        @description: start webrtc connection (invoked when user click on launch instance)
    */
    const callback_InitializeWebRTCConnection = (msg, payload) => {
        // obtain websocket connection from global map
        let ws = terminalWsMap.get(payload.TerminalKey)

        // send start streamming notification to scheduler
        let reqPacket = JSON.stringify({
            packet_type: "start_streamming",
            data: JSON.stringify({ 
                instance_id: payload.StateTerminals.terminalsMap[payload.TerminalKey].instanceSchedulerID,
            }),
        })
    }

    /*
        @function: registerRecvCallback
        @description: register recv callback functions of given websocket
    */
    const registerRecvCallback = () => {
        // user interaction callback
        PubSub.subscribe('config_websocket_state', callback_configWebsocketState)
        PubSub.subscribe('init_webrtc_connection', callback_InitializeWebRTCConnection)

        // websocket callback
        PubSub.subscribe('register_terminal_websocket', callback_registerTerminalWebsocket)
        PubSub.subscribe('delete_terminal_websocket', callback_deleteTerminalWebsocket)
        PubSub.subscribe('notify_ice_server', callback_recvNotifyIceServer)
        PubSub.subscribe('state_failed_provider_schedule', callback_stateFailedProviderSchedule)
        PubSub.subscribe('state_provider_scheduled', callback_stateProviderScheduled)
        PubSub.subscribe('state_failed_select_storage', callback_stateFailedStorageSchedule)
        PubSub.subscribe('state_selected_storage', callback_stateStorageScheduled)
        PubSub.subscribe('state_run_instance', callback_stateInstanceRunning)
        PubSub.subscribe('state_failed_run_instance', callback_stateFailedInstanceRunning)
    }

    useEffect( () => {
        registerRecvCallback()
        console.log("register ws receive callback!")
    },[])

    return <div />
}

export default WebsocketCallback
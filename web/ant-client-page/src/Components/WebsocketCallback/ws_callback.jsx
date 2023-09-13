import React, {useEffect, useState} from "react";
import PubSub from 'pubsub-js';
import GetTimestamp from '../../Utils/get_timestamp';
import { actions as TerminalActions } from '../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from "react-router-dom";
import { TERMINAL_STEP_SCHEDULE_COMPUTE_NODE, TERMINAL_STEP_CONFIG_INSTANCE, TERMINAL_STEP_SCHEDULE_STORAGE_NODE, TERMINAL_STEP_PREPARE_INSTANCE, TERMINAL_STEP_RUN_INSTANCE } from "../../Containers/UserPage/terminals";

const WebsocketCallback = (props) => {
    const dispatch = useDispatch()
    
    const navigate = useNavigate();

    const { terminalWsMap, setTerminalWsMap, terminalRtcPeerMap, setTerminalRtcPeerMap, terminalDynamicState, setTerminalDynamicState } = props;

    /*
        @callback: callback_registerTerminalDynamicState
        @description: 
            callback function for registering new "terminal <-> websocket" mapping relationship
    */
    const callback_registerTerminalDynamicState = (msg, payload) => {
        setTerminalDynamicState(terminalDynamicState.set(payload.TerminalKey, {
            instanceSchedulerID: "",
            iceServers: [],
            clientID: "",
        }));
    }

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
                    @ case: offer_sdp
                    @ description: WebRTC offer SDP from provider
                */
                case "offer_sdp":
                    PubSub.publish('offer_sdp', { 
                        Socket: ws,
                        TerminalKey: `${terminalKey}`,
                        WSPacket: wsPacket,
                        StateTerminals: stateTerminals,
                    });
                    break;

                /*
                    @ case: provider_ice_candidate
                    @ description: WebRTC ICE candidate from provider
                */
                case "provider_ice_candidate":
                    PubSub.publish('provider_ice_candidate', { 
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
        let instanceDynamicState = terminalDynamicState.get(`${payload.TerminalKey}`)

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

        // dispatch(TerminalActions.updateTerminal({
        //     "type": "APPEND_LOG_CONTENT",
        //     "terminal_key": `${payload.TerminalKey}`,
        //     "log_priority": "SUCCESS",
        //     "log_time": GetTimestamp(),
        //     "log_content": `User service has been successfully decomposed into three microservices: (1) normal computing power (2) app files for winmine (3) image layers for dcwine`,
        // }))

        // store iceserver
        instanceDynamicState.iceServers = payload.WSPacket.data.iceservers

        // store client id from scheduler
        instanceDynamicState.clientID = payload.WSPacket.data.client_id

        // save updated terminal dynamic state
        setTerminalDynamicState(terminalDynamicState.set(payload.TerminalKey, instanceDynamicState)) 

        // send metadata of selected application to scheduler
        let reqWSPacket = JSON.stringify({
            packet_type: "select_stream_application",
            data: JSON.stringify({ 
                application_id: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].applicationMeta.currentSelectedApplication.id}`,
                screen_height: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].screenHeight}`,
                screen_width: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].screenWidth}`,
                application_fps: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].currentFPS}`,
                vcodec: `${payload.StateTerminals.terminalsMap[payload.TerminalKey].vCodec}`,
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
            "log_content": `scheduler has found one provider to serve, provider id: ${payload.WSPacket.data.provider_id}, ip: ${payload.WSPacket.data.provider_ip}, is powerfull: ${payload.WSPacket.data.provider_is_powerful}`,
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
            "log_content": `provider has found proper depository worker node: address ${payload.WSPacket.data.depository_address}, is fast-speed: ${payload.WSPacket.data.depository_is_powerful}`,
        }))

        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "SUCCESS",
            "log_time": GetTimestamp(),
            "log_content": `provider has found proper filestore worker node: address ${payload.WSPacket.data.filestore_address}, is fast-speed: ${payload.WSPacket.data.filestore_is_powerful}`,
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

        // update instance id from scheduler
        let instanceDynamicState = terminalDynamicState.get(`${payload.TerminalKey}`)
        instanceDynamicState.instanceSchedulerID = payload.WSPacket.data.stream_instance_id

        // save updated terminal dynamic state
        setTerminalDynamicState(terminalDynamicState.set(payload.TerminalKey, instanceDynamicState))

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
        @callback: callback_registerTerminalRTCPeer
        @description: 
            callback function for registering new "terminal <-> RTCPeer" mapping relationship
    */
    const callback_registerTerminalRTCPeer = (msg, payload) => {
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, {
            PeerConnection: null,
            inputChannel: null,
            mediaStream: null,
            candidates: Array(),
            isAnswered: false,
            isFlushing: false,
            connected: false,
            inputReady: false
        }))
    }

    /*
        @callback: callback_InitializeWebRTCPeer
        @description: start webrtc connection (invoked when user click on launch instance)
    */
    const callback_InitializeWebRTCPeer = (msg, payload) => {
        // obtain RtcPeer object from global map
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)
        let instanceDynamicState = terminalDynamicState.get(`${payload.TerminalKey}`)

        // create new WebRTC peer connection
        let connection = new RTCPeerConnection({
            iceServers: JSON.parse(instanceDynamicState.iceServers)
        })
        RtcPeer.PeerConnection = connection

        // create new media stream
        // let mediaStream = new MediaStream()
        // RtcPeer.mediaStream = mediaStream

        /*
            @callback: ondatachannel
            @description: 
                register callbacks for the input channel
        */
        connection.ondatachannel = (e) => {
            RtcPeer.inputChannel = e.channel
            /*
                @callback: onopen
                @description: 
                    invoked when input channel is opened
            */
            RtcPeer.inputChannel.onopen = () => {
                // append log
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "INFO",
                    "log_time": GetTimestamp(),
                    "log_content": `Input channel attached to WebRTC connection is opened`,
                }))
                RtcPeer.inputReady = true
            }
            
            /*
                @callback: onopen
                @description: 
                    invoked when input channel is closed
            */
            RtcPeer.inputChannel.onclose = () => {
                // append log
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "WARN",
                    "log_time": GetTimestamp(),
                    "log_content": `Input channel attached to WebRTC connection is closed`,
                }))
                RtcPeer.inputReady = false
           }
        }

        /*
            @callback: oniceconnectionstatechange
            @description: 
                todo
        */
        connection.oniceconnectionstatechange = (e) => {
            PubSub.publish('webrtc_oniceconnectionstatechange', { 
                TerminalKey: `${payload.TerminalKey}`,
                Event: e,
            });
        }

        /*
            @callback: onicegatheringstatechange
            @description: 
                todo
        */
        connection.onicegatheringstatechange = (e) => {
            PubSub.publish('webrtc_onicegatheringstatechange', { 
                TerminalKey: `${payload.TerminalKey}`,
                Event: e,
            });
        }

        /*
            @callback: onicecandidate
            @description: 
                todo
        */
        connection.onicecandidate = (e) => {
            PubSub.publish('webrtc_onicecandidate', { 
                TerminalKey: `${payload.TerminalKey}`,
                Event: e,
            });
        }

        /*
            @callback: ontrack
            @description: 
                todo
        */
        connection.ontrack = (e) => {
            PubSub.publish('webrtc_ontrack', { 
                TerminalKey: `${payload.TerminalKey}`,
                Streams: e.streams,
            });
        }

        // save updated RtcPeer
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))

        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `created WebRTC peer connection for current terminal`,
        }))

        // obtain websocket connection from global map
        let ws = terminalWsMap.get(payload.TerminalKey)

        // send start streamming notification to scheduler
        let reqPacket = JSON.stringify({
            packet_type: "start_streaming",
            data: JSON.stringify({ 
                instance_id: terminalDynamicState.get(payload.TerminalKey).instanceSchedulerID
            }),
        })

        // send to scheudler
        ws.send(reqPacket)
    }

    /*
        @callback: callback_ProviderOfferSDP
        @description: 
            invoked while receiving provider offer SDP
    */
    const callback_ProviderOfferSDP = async (msg, payload) => {
        // obtain RTC Peer object
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)

        // parse offer SDP
        const offer = new RTCSessionDescription(JSON.parse(atob(payload.WSPacket.data.offer_sdp)))
        
        // set remote description
        await RtcPeer.PeerConnection.setRemoteDescription(offer)

        // create answer sdp
        const answer = await RtcPeer.PeerConnection.createAnswer()
        answer.sdp = answer.sdp.replace(
            /(a=fmtp:111 .*)/g,
            "$1;stereo=1;sprop-stereo=1"
        )
        await RtcPeer.PeerConnection.setLocalDescription(answer);
        RtcPeer.PeerConnection.isAnswered = true
        
        // save updated RtcPeer
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))
            
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `receive offer SDP from provider, set as remote description`,
        }))

        // obtain websocket connection from global map
        let ws = terminalWsMap.get(payload.TerminalKey)

        // send start streamming notification to scheduler
        let reqPacket = JSON.stringify({
            packet_type: "answer_sdp",
            data: JSON.stringify({ 
                answer_sdp: btoa(JSON.stringify(answer)),
                instance_id: terminalDynamicState.get(payload.TerminalKey).instanceSchedulerID,
            }),
        })
        
        // send to scheudler
        ws.send(reqPacket)
        
        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `send answer SDP back to provider`,
        }))
    }

    /*
        @callback: callback_ProviderICECandidate
        @description: 
            invoked while receiving provider ice candidate
    */
    const callback_ProviderICECandidate = async (msg, payload) => {
        // obtain RTC Peer object
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)

        // decode
        let candidate_decode = atob(payload.WSPacket.data.provider_ice_candidate);

        if(candidate_decode !== null && candidate_decode !== ""){
            if(RtcPeer.isAnswered === false){
                // add ice candidate
                let candidate = new RTCIceCandidate(JSON.parse(candidate_decode));
                RtcPeer.PeerConnection.addIceCandidate(candidate);

                // save updated RtcPeer
                setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))

                // append log
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "INFO",
                    "log_time": GetTimestamp(),
                    "log_content": `add ice candidate of remote provider: ${JSON.parse(atob(payload.WSPacket.data.provider_ice_candidate)).candidate}`,
                }))
            } else {
                // append log
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "WARN",
                    "log_time": GetTimestamp(),
                    "log_content": `ignore ice candidate of remote provider: ${JSON.parse(atob(payload.WSPacket.data.provider_ice_candidate)).candidate}`,
                }))
            }
        }
    }

    /*
        @callback: callback_FlushICECandidate
        @description: 
            invoked while receiving empty provider ice candidate
    */
    const callback_FlushICECandidate = async (msg, payload) => {
        // obtain RTC Peer object
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)
        if (RtcPeer.isFlushing || !RtcPeer.isAnswered)
            return
        
        RtcPeer.isFlushing = true
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))

        RtcPeer.candidates.forEach(data => {
            let d = atob(data);
            let candidate = new RTCIceCandidate(JSON.parse(d));
            RtcPeer.PeerConnection.addIceCandidate(candidate);
        });

        RtcPeer.isFlushing = false

        // save updated RtcPeer
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))

        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `flush existed candidates`,
        }))
    }

    /*
        @callback: callback_ICEConnectionStateChange
        @description: 
            invoked while state of WebRTC ICE connection changed
    */
    const callback_ICEConnectionStateChange = async (msg, payload) => {
        // obtain RTC Peer object
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)
        
        switch(RtcPeer.PeerConnection.iceConnectionState){
            case "connected":
                RtcPeer.connected = true
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "INFO",
                    "log_time": GetTimestamp(),
                    "log_content": 'ICE connectition state changed: connected',
                }))

                // navigate to stream page
                navigate({
                    pathname: '/stream',
                    search: `?key=${payload.TerminalKey}`
                })

                break
            
            case "disconnected":
                RtcPeer.connected = false
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "WARN",
                    "log_time": GetTimestamp(),
                    "log_content": 'ICE connectition state changed: disconnected',
                }))
                break
            
            case "failed":
                RtcPeer.connected = false
                dispatch(TerminalActions.updateTerminal({
                    "type": "APPEND_LOG_CONTENT",
                    "terminal_key": `${payload.TerminalKey}`,
                    "log_priority": "ERROR",
                    "log_time": GetTimestamp(),
                    "log_content": 'ICE connectition state changed: failed, restart',
                }))
                const offer = await RtcPeer.PeerConnection.createOffer({ iceRestart: true })
                RtcPeer.PeerConnection.setLocalDescription(offer)
                break
        }

        // save updated RtcPeer
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))
    }

    /*
        @state: intervalForGathering
        @description: interval for ice candidate gathering
    */
    const [intervalForGathering, setIntervalForGathering] = useState(null)

    /*
        @callback: callback_ICEGatheringStateChange
        @description: 
            invoked while state of WebRTC ICE gathering state change
    */
    const callback_ICEGatheringStateChange = (msg, payload) => {
        const ICE_GATHERING_TIMEOUT = 2000
        switch(payload.Event.target.iceGatheringState){
            case "gathering":
                let interval = setTimeout(() => {
                    // append log
                    dispatch(TerminalActions.updateTerminal({
                        "type": "APPEND_LOG_CONTENT",
                        "terminal_key": `${payload.TerminalKey}`,
                        "log_priority": "INFO",
                        "log_time": GetTimestamp(),
                        "log_content": `ICE gathering state exceed maximum duration, timeout`,
                    }))
                }, ICE_GATHERING_TIMEOUT)
                setIntervalForGathering(interval)
                break
            
            case "complete":
                if(intervalForGathering){
                    clearTimeout(intervalForGathering)
                }
                break
        }
    }

    /*
        @callback: callback_onICECandidate
        @description: 
            invoked while notification of ice candidates from STUN server
    */
    const callback_onICECandidate = (msg, payload) => {
        // obtain websocket connection from global map
        let ws = terminalWsMap.get(payload.TerminalKey)

        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `received ICE candidate infomation from ICE server: ${JSON.stringify(payload.Event.candidate)}`,
        }))

        // send start streamming notification to scheduler
        if(payload.Event.candidate !== null){
            let reqPacket = JSON.stringify({
                packet_type: "consumer_ice_candidate",
                data: JSON.stringify({ 
                    instance_id: terminalDynamicState.get(payload.TerminalKey).instanceSchedulerID,
                    consumer_ice_candidate: btoa(JSON.stringify(payload.Event.candidate))
                }),
            })

            // send to scheudler
            ws.send(reqPacket)
        }
    }

    /*
        @callback: callback_onTrack
        @description: 
            todo
    */
    const callback_onTrack = (msg, payload) => {
        let RtcPeer = terminalRtcPeerMap.get(payload.TerminalKey)

        // add track
        // RtcPeer.mediaStream.addTrack(payload.Track);
        RtcPeer.mediaStream = payload.Streams

        // save updated RtcPeer
        setTerminalRtcPeerMap(terminalRtcPeerMap.set(payload.TerminalKey, RtcPeer))

        // append log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${payload.TerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `add new track to media stream`,
        }))
    }

    /*
        @function: registerRecvCallback
        @description: register recv callback functions of given websocket
    */
    const registerRecvCallback = () => {
        // user interaction callback (websocket)
        PubSub.subscribe('config_websocket_state', callback_configWebsocketState)
        PubSub.subscribe('register_terminal_websocket', callback_registerTerminalWebsocket)
        PubSub.subscribe('delete_terminal_websocket', callback_deleteTerminalWebsocket)
        PubSub.subscribe('register_terminal_dynamic_state', callback_registerTerminalDynamicState)

        // user interaction callback (webRTC)
        PubSub.subscribe('init_webrtc_connection', callback_InitializeWebRTCPeer)
        PubSub.subscribe('register_terminal_rtc_connection', callback_registerTerminalRTCPeer)
        
        // websocket callback
        PubSub.subscribe('notify_ice_server', callback_recvNotifyIceServer)
        PubSub.subscribe('state_failed_provider_schedule', callback_stateFailedProviderSchedule)
        PubSub.subscribe('state_provider_scheduled', callback_stateProviderScheduled)
        PubSub.subscribe('state_failed_select_storage', callback_stateFailedStorageSchedule)
        PubSub.subscribe('state_selected_storage', callback_stateStorageScheduled)
        PubSub.subscribe('state_run_instance', callback_stateInstanceRunning)
        PubSub.subscribe('state_failed_run_instance', callback_stateFailedInstanceRunning)
        PubSub.subscribe('offer_sdp', callback_ProviderOfferSDP)
        PubSub.subscribe('provider_ice_candidate', callback_ProviderICECandidate)

        // webRTC callback
        PubSub.subscribe('webrtc_oniceconnectionstatechange', callback_ICEConnectionStateChange)
        PubSub.subscribe('webrtc_onicegatheringstatechange', callback_ICEGatheringStateChange)
        PubSub.subscribe('webrtc_onicecandidate', callback_onICECandidate)
        PubSub.subscribe('webrtc_ontrack', callback_onTrack)
        PubSub.subscribe('flush_ice_candidate', callback_FlushICECandidate)        
    }

    useEffect( () => {
        registerRecvCallback()
        console.log("register ws receive callback!")
    },[])

    return null
}

export default WebsocketCallback
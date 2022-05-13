import React from 'react';
import styled from 'styled-components';
import PubSub from 'pubsub-js';
import Button from '@mui/material/Button';
import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import PauseCircleIcon from '@mui/icons-material/PauseCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import LaunchIcon from '@mui/icons-material/Launch';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../../Data/Reducers/snackBarReducer'
import GetTimestamp from '../../../Utils/get_timestamp';
import { TabIndex_Dashboard_Application, TabIndex_Dashboard_LogViewer, TabIndex_Dashboard_Terminal } from '../terminalConfig';
import { TERMINAL_STEP_CONFIG_INSTANCE, TERMINAL_STEP_CONNECT_TO_SCHEDULER, TERMINAL_STEP_PREPARE_INSTANCE, TERMINAL_STEP_RUN_INSTANCE, TERMINAL_STEP_SCHEDULE_COMPUTE_NODE, TERMINAL_STEP_SCHEDULE_STORAGE_NODE } from '../terminals';

const ControlPanelBtnGroupContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-evenly;
    border-radius: 10px;
    background-color: #e8e8e8;
    padding-top: 30px;
`

const ButtonItem = styled.div`
    display: flex;
    align-items: center;
    margin-bottom: 30px;
`

const ButtonItemDesp = styled.p`
    font-size: 10px;
    color: #8c8c8c;
    margin: 3px 0px;
`


const TerminalControlPanel = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    const CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta

    // get global states of websocket api
    const StateAPI = useSelector(state => state.api.StateAPI)
    const StateInfo = useSelector(state => state.info.StateInfo)

    /*
        @function: handleTerminalCreate
        @description:
            handle terminal creation
    */
    const handleTerminalCreate = (event) => {
        // check whether user have selected application
        if(!CurrentApplicationMeta.selectedApplicationConfirmed){
            // show snackbar
            dispatch(SnackBarActions.showSnackBar(`Please choose your desired application before you create this terminal`))
            
            // jump to application panel
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_DASHBOARD_ID",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "dashboard_id": TabIndex_Dashboard_Application
            }))
            
            return
        }

        // check whether user have configed terminal
        if(!CurrentSelectedTerminal.terminalConfigConfirm){
            // show snackbar
            dispatch(SnackBarActions.showSnackBar(`Please config related terminal metadata before you create this terminal`))
            
            // jump to application panel
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_DASHBOARD_ID",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "dashboard_id": TabIndex_Dashboard_Terminal
            }))

            return
        }

        // confirm websocket connection started
        dispatch(TerminalActions.updateTerminal({
            "type": "CONFIRM_WS_CONNECTION_STARTED",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))

        // jump to the log panel
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_DASHBOARD_ID",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "dashboard_id": TabIndex_Dashboard_LogViewer
        }))

        // record current selected terminal, used in callbacks
        let currentSelectedTerminalKey = StateTerminals.currentSelected
        let currentSelectedTerminalName = CurrentSelectedTerminal.name

        // create websocket
        const ws = new WebSocket(`${StateAPI.WSProtocol}://${StateAPI.WSHostAddr}:${StateAPI.WSPort}${StateAPI.WSBaseURL}/${StateAPI.WSAPI.WebSocketConnect}?type=${StateInfo.ClientType}`)
        
        // append terminal log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${currentSelectedTerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": `try to connect to scheduler: ${StateAPI.WSProtocol}://${StateAPI.WSHostAddr}:${StateAPI.WSPort}${StateAPI.WSBaseURL}/${StateAPI.WSAPI.WebSocketConnect}?type=${StateInfo.ClientType}`,
        }))

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${currentSelectedTerminalKey}`,
            "current_step_index": TERMINAL_STEP_CONNECT_TO_SCHEDULER
        }))

        // register websocket
        PubSub.publish('register_terminal_websocket', {
            Websocket: ws,
            TerminalKey: currentSelectedTerminalKey
        });

        // register terminal dynamic state
        PubSub.publish('register_terminal_dynamic_state', {
            TerminalKey: currentSelectedTerminalKey
        });

        // config the behaviors under different states of newly created websocket
        PubSub.publish('config_websocket_state', {
            Websocket: ws,
            TerminalKey: currentSelectedTerminalKey,
            TerminalName: currentSelectedTerminalName,
            StateTerminals: StateTerminals
        });
    }

    /*
        @function: handleTerminalCancelConnection
        @description:
            handle cancel connection to scheduler
    */
    const handleTerminalCancelConnection = (event) => {
        // record current selected terminal, used in callbacks
        const currentSelectedTerminalKey = StateTerminals.currentSelected

        // append terminal log
        dispatch(TerminalActions.updateTerminal({
            "type": "APPEND_LOG_CONTENT",
            "terminal_key": `${currentSelectedTerminalKey}`,
            "log_priority": "INFO",
            "log_time": GetTimestamp(),
            "log_content": "user shutdown websocket connection",
        }))
        
        // close websocket
        PubSub.publish('delete_terminal_websocket', {
            TerminalKey: currentSelectedTerminalKey
        });

        // change current step
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_STEP",
            "terminal_key": `${currentSelectedTerminalKey}`,
            "current_step_index": TERMINAL_STEP_CONFIG_INSTANCE
        }))

        // unconfirm websocket connection started
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_WS_CONNECTION_STARTED",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))
    }

    /*
        @function: handleTerminalDelete
        @description:
            handle delete terminal
    */
    const handleTerminalDelete = (event, newValue) => {
        // close connection to scheduler
        handleTerminalCancelConnection()

        // delete terminal
        dispatch(TerminalActions.deleteTerminal(StateTerminals.currentSelected))
    }

    /*
        @function: handleLaunchInstance
        @description:
            handle launch instance
    */
    const handleLaunchInstance = (event, newValue) => {
        // resgiter rtc connection
        PubSub.publish('register_terminal_rtc_connection', { 
            TerminalKey: `${StateTerminals.currentSelected}`,
        });

        // initilize webrtc connection
        PubSub.publish('init_webrtc_connection', { 
            TerminalKey: `${StateTerminals.currentSelected}`,
            StateTerminals: StateTerminals,
        });
    }

    return <ControlPanelBtnGroupContainer>
        {/* Start || CancelConnect */}
        <ButtonItem>
        {
            CurrentSelectedTerminal.currentStepIndex === TERMINAL_STEP_CONFIG_INSTANCE &&
            <div align="center">
                <Button 
                    variant="contained"
                    color="success"
                    startIcon={<PlayCircleFilledIcon />}
                    onClick={handleTerminalCreate}
                >
                    Create
                </Button>
                <ButtonItemDesp>Create Instance</ButtonItemDesp>
            </div>
        }

        {
            (CurrentSelectedTerminal.currentStepIndex !== TERMINAL_STEP_CONFIG_INSTANCE) &&
            <div align="center">
                <Button 
                    variant="contained"
                    color="warning"
                    startIcon={<PauseCircleIcon />}
                    onClick={handleTerminalCancelConnection}
                >
                    Cancel
                </Button>
                <ButtonItemDesp>Cancel Connection</ButtonItemDesp>
            </div>
        }
        </ButtonItem>

        {/* Launch */}
        <ButtonItem>
        <div align="center">
            <Button 
                variant="contained"
                startIcon={<LaunchIcon />}
                disabled={CurrentSelectedTerminal.currentStepIndex !== TERMINAL_STEP_RUN_INSTANCE}
                onClick={handleLaunchInstance}
            >
                Launch
            </Button>
            <ButtonItemDesp>Launch Instance</ButtonItemDesp>
        </div>
        </ButtonItem>

        {/* Delete */}
        <ButtonItem>
        <div align="center">
            <Button 
                variant="contained"
                startIcon={<DeleteIcon />}
                disabled={CurrentSelectedTerminal.currentStepIndex !== TERMINAL_STEP_CONFIG_INSTANCE}
                onClick={handleTerminalDelete}
            >
                Delete
            </Button>
            <ButtonItemDesp>Delete Terminal</ButtonItemDesp>
        </div>
        </ButtonItem>
    </ControlPanelBtnGroupContainer>
}

export default TerminalControlPanel
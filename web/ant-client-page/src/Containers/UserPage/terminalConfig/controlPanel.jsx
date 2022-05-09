import React from 'react';
import styled from 'styled-components';
import PubSub from 'pubsub-js';
import Button from '@mui/material/Button';
import PlayCircleFilledIcon from '@mui/icons-material/PlayCircleFilled';
import PauseCircleIcon from '@mui/icons-material/PauseCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../../Data/Reducers/snackBarReducer'
import GetTimestamp from '../../../Utils/get_timestamp';
import { TabIndex_Dashboard_Application, TabIndex_Dashboard_LogViewer } from '../terminalConfig';

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
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    let CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta

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
            "current_step_index": 1
        }))

        // register websocket
        PubSub.publish('register_terminal_websocket', {
            Websocket: ws,
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
            "current_step_index": 0
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

    return <ControlPanelBtnGroupContainer>
        {/* Start || CancelConnect */}
        <ButtonItem>
        {
            CurrentSelectedTerminal.currentStepIndex === 0 &&
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
            (CurrentSelectedTerminal.currentStepIndex === 1 || CurrentSelectedTerminal.currentStepIndex === 2) &&
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

        {/* Delete */}
        <ButtonItem>
        <div align="center">
            <Button 
                variant="contained"
                startIcon={<DeleteIcon />}
                disabled={CurrentSelectedTerminal.currentStepIndex === 1 || CurrentSelectedTerminal.currentStepIndex === 4}
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
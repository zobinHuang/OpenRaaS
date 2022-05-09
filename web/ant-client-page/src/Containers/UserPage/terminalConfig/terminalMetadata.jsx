import React from 'react';
import styled from 'styled-components';
import TextField from '@mui/material/TextField';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../Data/Reducers/terminalReducer';

const TerminalMetadata = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    let StateTerminals = useSelector(state => state.terminal.StateTerminals)
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]

     /*
        @function: handleTerminalNameUpdate
        @description:
            handle update terminal name
    */
    const handleTerminalNameUpdate = (event, newValue) => {
        dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_NAME",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "name": event.target.value
        }))
    }

    return <TextField
        value={CurrentSelectedTerminal.name}
        disabled={CurrentSelectedTerminal.currentStepIndex !== 0}
        onChange={handleTerminalNameUpdate}
        style={{zIndex: "0"}}
        id="terminal-name" 
        label="Terminal Name" 
        variant="outlined"
        size="small"
        fullWidth
    />
}

export default TerminalMetadata
import React, { useEffect, useState } from 'react';
import IconButton from '@mui/material/IconButton';
import AddIcon from '@mui/icons-material/AddCircle';
import Slide from '@mui/material/Slide';
import styled from 'styled-components';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Tooltip from '@mui/material/Tooltip';
import BuildCircleIcon from '@mui/icons-material/BuildCircle';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import PendingIcon from '@mui/icons-material/Pending';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import FreeBreakfastIcon from '@mui/icons-material/FreeBreakfast';
import { useSelector, useDispatch } from 'react-redux'
import { actions as TerminalActions } from '../../Data/Reducers/terminalReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import TerminalConfigList from './terminalConfig';

/* Container: Header */
const TerminalsPageContainer = styled.div`
    width: 95%;
    margin: 0px;
    padding: 0px;
    display: flex;
    box-shadow: 5px 5px 5px #a0a0a0;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-color: #e3f9ff;
`

const TabContainer = styled.div`
    width: 100%;
    height: 12%;
    margin: 0px;
    display: flex;
    align-items: center;
    justify-content: center;
`

const TabGroupContainer = styled.div`
    width: 95%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const TabButtonContainer = styled.div`
    width: 5%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const TerminalConfigContainer = styled.div`
    width: 100%;
    background-color: #ffffff;
    overflow-y: scroll;
    overflow-x: hidden;
`

const EmptyTerminalContainer = styled.div`
    width: 100%;
    min-height: calc(60vh);
    background-color: #f2f2f2;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const EmptyTerminalPrompt = styled.h2`
    color: #c4c4c4;
`

export const TERMINAL_STEP_CONFIG_INSTANCE = 0
export const TERMINAL_STEP_CONNECT_TO_SCHEDULER = 1
export const TERMINAL_STEP_SCHEDULE_COMPUTE_NODE = 2
export const TERMINAL_STEP_SCHEDULE_STORAGE_NODE = 3
export const TERMINAL_STEP_PREPARE_INSTANCE = 4
export const TERMINAL_STEP_RUN_INSTANCE = 5

const TerminalsPage = (props) => {    

    const dispatch = useDispatch()

    // Get global state
    const StateTerminals = useSelector( state => state.terminal.StateTerminals )
    const StateTerminalMap = useSelector(state => state.terminal.StateTerminals.terminalsMap)

    // Handle tab change
    const handleChangeSelectedTerminal = (event, newValue) => {
        event.preventDefault()
        dispatch(TerminalActions.changeSelectedTab(newValue))
    };

    // Handle add terminal
    const handleAddTerminal = (event) => {
        event.preventDefault()

        // check terminal amount (not larger than 5)
        if(Object.getOwnPropertyNames(StateTerminalMap).length >= 5){
            dispatch(SnackBarActions.showSnackBar("One can only create 5 terminals at most"))
            return
        }

        // add new terminal in store
        dispatch(TerminalActions.addTerminal())
    }

    return (
        <Slide 
            direction="down" 
            in={true}
            mountOnEnter 
            unmountOnExit
        >
            <TerminalsPageContainer
                id="terminal-page-container"
            >
                {/* Tab Area */}
                <TabContainer>
                    {/* Tab Group */}
                    <TabGroupContainer>
                    <Tabs
                        value={StateTerminals.currentSelectedIndex}
                        onChange={handleChangeSelectedTerminal}
                        variant="scrollable"
                        scrollButtons="auto"
                        aria-label="scrollable terminals tab"
                    >
                        {Object.getOwnPropertyNames(StateTerminalMap).length > 0 ?
                            Object.keys(StateTerminalMap).map(
                                (terminalIndex, index) => {
                                    let terminal = StateTerminalMap[terminalIndex]
                                    switch(terminal.currentStepIndex){
                                        case TERMINAL_STEP_CONFIG_INSTANCE:
                                            return <Tab 
                                                icon={<BuildCircleIcon />}
                                                key={terminalIndex}
                                                iconPosition="start"
                                                label={terminal.name}
                                                style={{textTransform: 'none'}}
                                            />
                                        case TERMINAL_STEP_CONNECT_TO_SCHEDULER:
                                            return <Tab 
                                                icon={<ArrowCircleUpIcon />}
                                                key={terminalIndex}
                                                iconPosition="start"
                                                label={terminal.name}
                                                style={{textTransform: 'none'}}
                                            />
                                        case TERMINAL_STEP_RUN_INSTANCE:
                                            return <Tab 
                                                icon={<CheckCircleIcon />}
                                                key={terminalIndex}
                                                iconPosition="start"
                                                label={terminal.name}
                                                style={{textTransform: 'none'}}
                                            />
                                        default:
                                            return <Tab 
                                                icon={<PendingIcon />}
                                                key={terminalIndex}
                                                iconPosition="start"
                                                label={terminal.name}
                                                style={{textTransform: 'none'}}
                                            />
                                    }
                                }
                            )
                        : <Tab 
                            label="Empty Terminal List"
                            style={{textTransform: 'none'}}
                        />}
                    </Tabs>
                    </TabGroupContainer>

                    {/* Add Button */}
                    <TabButtonContainer>
                    <Tooltip title="Add New Terminal">
                    <IconButton 
                        aria-label="add-terminal"
                        onClick={handleAddTerminal}
                    >
                        <AddIcon fontSize="large" />
                    </IconButton>
                    </Tooltip>
                    </TabButtonContainer>

                </TabContainer>
                
                {
                    Object.getOwnPropertyNames(StateTerminalMap).length > 0 ? 
                    // Terminal Configuration List
                    <TerminalConfigContainer>
                        <TerminalConfigList />
                    </TerminalConfigContainer> : 

                    // Empty Terminal Prompt
                    <EmptyTerminalContainer>
                        <FreeBreakfastIcon color="disabled" sx={{ fontSize: 300, margin: 0 }} />
                        <EmptyTerminalPrompt>No Terminal is Running</EmptyTerminalPrompt>
                    </EmptyTerminalContainer>
                }
            </TerminalsPageContainer>
        </Slide>
    )
}

export default TerminalsPage
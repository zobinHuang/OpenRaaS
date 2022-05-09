import React from 'react';
import styled from 'styled-components';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import { useDispatch, useSelector } from 'react-redux';
import LogoDevIcon from '@mui/icons-material/LogoDev';
import CloudCircleIcon from '@mui/icons-material/CloudCircle';
import TerminalIcon from '@mui/icons-material/Terminal';
import SpeedIcon from '@mui/icons-material/Speed';
import { actions as TerminalActions } from '../../Data/Reducers/terminalReducer';
import RollStepper from '../../Components/RollStepper/rollStepper';
import LogViewer from '../../Components/LogTerminal/logTerminal';
import TerminalControlPanel from './terminalConfig/controlPanel';
import TerminalMetadata from './terminalConfig/terminalMetadata';
import ApplicationMetadata from './terminalConfig/applicationMetadata';
import Badge from '@mui/material/Badge';

const TerminalConfigContainer = styled.div`
    width: 100%;
    height: 100%;
    margin: 0px;
    padding: 20px 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
`

const TitleContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: center;
    width: 92%;
    margin-bottom: 5px;
    background-color: #0057a3;
    padding: 8px 20px;
    border-radius: 8px;
`

const Title = styled.h2`
    color: #ffffff;
    margin: 0px;
`

const StepperContainer = styled.div`
    width: 90%;
    height: 120px;
    margin-bottom: 50px;
    padding: 20px 30px;
    border: 1px solid #b0b0b0;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 2px 2px 2px #a0a0a0;
`

const DashboardContainer = styled.div`
    display: flex;
    flex-direction: row;
    width: 90%;
    min-height: 500px;
    padding: 20px 30px;
    margin-bottom: 20px;
    border: 1px solid #b0b0b0;
    border-radius: 10px;
    box-shadow: 2px 2px 2px #a0a0a0;
`

const TerminalContrlPanelContainer = styled.div`
    display: flex;
    width: 90%;
    align-items: center;
    justify-content: center;
    padding: 0px 30px;
    background-color: #e8e8e8;
    margin-top: 5px;
    margin-bottom: 20px;
`

const TabContainer = styled.div`
    width: 10%;
    display: flex;
    // align-items: center;
    justify-content: center;
`

const DetailContainer = styled.div`
    width: 90%;
`

/*
    @enum: TabIndex
    @description:
        tab index of different dashboard
*/
export const TabIndex_Dashboard_Application = 0
export const TabIndex_Dashboard_Terminal = 1
export const TabIndex_Dashboard_LogViewer = 3

/*
    @component: TerminalConfigList
    @description:
        terminal configuration list
*/
const TerminalConfigList = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]

    // RollStepperConfig
    const RollStepperConfig = {
        steps: CurrentSelectedTerminal.steps,
        currentStepIndex: CurrentSelectedTerminal.currentStepIndex
    }

    /*
        @function: handleDashboardTabChange
        @description:
            handle update selected dashboard tab index
    */
    const handleDashboardTabChange = (event, newValue) => {
        dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_DASHBOARD_ID",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "dashboard_id": newValue
        }))

        // clear unread log count
        if(newValue === TabIndex_Dashboard_LogViewer){
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_UNREAD_LOG_COUNT",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "unread_log_count": 0
            }))
        }
    }

    return (
        <TerminalConfigContainer >
            {/* Stepper */}
            <TitleContainer><Title>Process</Title></TitleContainer>
            <StepperContainer>
                <RollStepper
                    RollStepperConfig={RollStepperConfig}
                />
            </StepperContainer>
            
            {/* Dashboard */}
            <TitleContainer><Title>Dashboard</Title></TitleContainer>
            
            <TerminalContrlPanelContainer><TerminalControlPanel /></TerminalContrlPanelContainer>

            <DashboardContainer>
                <TabContainer>
                    <Tabs
                        orientation="vertical"
                        variant="scrollable"
                        value={CurrentSelectedTerminal.currentDashboardIndex}
                        onChange={handleDashboardTabChange}
                        aria-label="Vertical tabs example"
                        indicatorColor="secondary"
                        sx={{ borderRight: 1, borderColor: 'divider'}}
                    >
                        <Tab 
                            icon={<CloudCircleIcon />}
                            key={'dashboard-applictaion'}
                            iconPosition="start"
                            label={'App'}
                            style={{textTransform: 'none'}}
                        />
                        <Tab 
                            icon={<TerminalIcon />}
                            key={'dashboard-terminal'}
                            iconPosition="start"
                            label={'Terminal'}
                            style={{textTransform: 'none'}}
                        />
                        <Tab 
                            icon={<SpeedIcon />}
                            key={'dashboard-performance'}
                            iconPosition="start"
                            label={'Monitor'}
                            style={{textTransform: 'none'}}
                        />
                        <Tab 
                            icon={<LogoDevIcon />}
                            key={'dashboard-log'}
                            iconPosition="start"
                            label={<Badge 
                                badgeContent={CurrentSelectedTerminal.unreadLogCount} 
                                color={CurrentSelectedTerminal.unreadLogLevel}
                                sx={{'& .MuiBadge-badge': {
                                    right: -10,
                                    top: 12,
                                    border: `2px solid #ffffff`,
                                    padding: '0 4px',
                                  }}}
                            >
                                <p>Log</p>
                            </Badge>}
                            style={{textTransform: 'none'}}
                        />  
                    </Tabs>
                </TabContainer>
                <DetailContainer>
                    {
                        CurrentSelectedTerminal.currentDashboardIndex === TabIndex_Dashboard_Application && <ApplicationMetadata />
                    }
                    {
                        CurrentSelectedTerminal.currentDashboardIndex === TabIndex_Dashboard_Terminal && <TerminalMetadata />
                    }
                    {
                        CurrentSelectedTerminal.currentDashboardIndex === TabIndex_Dashboard_LogViewer && <LogViewer 
                            LogViewerConfiguration={{
                                entries: CurrentSelectedTerminal.logInfo,
                                height: "450px"
                            }}
                        />
                    }
                </DetailContainer>
                
            </DashboardContainer>
        </TerminalConfigContainer>
    )
}

export default TerminalConfigList
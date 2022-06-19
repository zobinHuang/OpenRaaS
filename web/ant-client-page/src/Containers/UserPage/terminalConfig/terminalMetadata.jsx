import React from 'react';
import styled from 'styled-components';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import Button from '@mui/material/Button';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../Data/Reducers/terminalReducer';

const TerminalMetadataContainer = styled.div`
    width: 100%;
    padding: 20px 20px;
    display: flex;
    flex-direction: column;
    margin-bottom: 20px;
`

const ConfigItemTitleContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    margin-top: 20px;
`

const ConfigItemTitle = styled.h3`
    margin: 0px;
    margin-bottom: 8px;
`

const TerminalNameConfigContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
`

const ResolutionConfigContainer = styled.div`
    width: 100%;
    border: 1px solid #969696;
    margin-bottom: 20px;
    padding: 20px 0px;
    display: flex;
    align-items: center;
    justify-content: space-evenly;
`

const ResolutionConfigEntry = styled.div`
    width: 20%;
    min-height: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
`

const ResolutionShowcase = styled.div`
    width: ${ ({Width}) => Width ? Width : "100%" };
    height: ${ ({Height}) => Height ? Height : "100%" };
    background-color: ${ ({Selected}) => Selected ? "#ff7474d3" : "#d2d2d2d4" };
    transition: 300ms;
    color: #000000;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: larger;
    font-weight: bolder;

    &:hover {
        box-shadow: 5px 5px 5px #a0a0a0;
        cursor: pointer;
    }
`

const ButtonGroupContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 20px 0px;
`


export const RESOLUTION_1280_720 = "resolution_1280_720"
export const RESOLUTION_640_480 = "resolution_640_480"
export const RESOLUTION_1024_600 = "resolution_1024_600"
export const RESOLUTION_1280_1024 = "resolution_1280_1024"

export const FPS_5 = 5
export const FPS_15 = 15
export const FPS_24 = 24
export const FPS_30 = 30
export const FPS_60 = 60
export const FPS_120 = 120

export const VCODEC_H264 = "h264"
export const VCODEC_VPX = "vpx"

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

    /*
        @function: handleSelect_640_480
        @description:
            handle update resolution to 640x480
    */
    const handleSelect_640_480 = (event) => {
        if(!CurrentSelectedTerminal.terminalConfigConfirm){
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_RESOLUTION",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "resolution": RESOLUTION_640_480,
                "height": 480,
                "width": 640
            }))
        }
    }

    /*
        @function: handleSelect_1024_600
        @description:
            handle update resolution to 1024x600
    */
    const handleSelect_1024_600 = (event) => {
        if(!CurrentSelectedTerminal.terminalConfigConfirm){
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_RESOLUTION",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "resolution": RESOLUTION_1024_600,
                "height": 600,
                "width": 1024
            }))
        }
    }

    /*
        @function: handleSelect_1280_1024
        @description:
            handle update resolution to 1280x1024
    */
    const handleSelect_1280_1024 = (event) => {
        if(!CurrentSelectedTerminal.terminalConfigConfirm){
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_RESOLUTION",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "resolution": RESOLUTION_1280_1024,
                "height": 1024,
                "width": 1280
            }))
        }
    }

    /*
        @function: handleSelect_1280_720
        @description:
            handle update resolution to 1280x720
    */
    const handleSelect_1280_720 = (event) => {
        if(!CurrentSelectedTerminal.terminalConfigConfirm){
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_TERMINAL_RESOLUTION",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "resolution": RESOLUTION_1280_720,
                "height": 720,
                "width": 1280
            }))
        }
    }

    /*
        @function: handleChangeFPS
        @description:
            handle update fps
    */
    const handleChangeFPS = (event) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_FPS",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "fps": event.target.value
        }))
    }

    /*
        @function: handleConfirmTerminalConfig
        @description:
            handle confirm terminal configuration
    */
    const handleConfirmTerminalConfig = (event) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "CONFIRM_TERMINAL_CONFIG",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))
    }

    /*
        @function: handleCancelTerminalConfig
        @description:
            handle cancel terminal configuration
    */
    const handleCancelTerminalConfig = (event) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_TERMINAL_CONFIG",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))
    }

    /*
        @function: handleChangeVCodec
        @description:
            handle change vcodec
    */
    const handleChangeVCodec = (event) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_TERMINAL_VCODEC",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "vcodec": event.target.value,
        }))
    }

    const fpsFormat = [
        {
            value: FPS_5,
            label: "5fps"
        },
        {
            value: FPS_15,
            label: "15fps"
        },
        {
            value: FPS_24,
            label: "24fps"
        },
        {
            value: FPS_30,
            label: "30fps"
        },
        {
            value: FPS_60,
            label: "60fps"
        },
        {
            value: FPS_120,
            label: "120fps"
        },
    ]

    const vcodecFormat = [
        {
            value: VCODEC_H264,
            label: "H.264"
        },
        {
            value: VCODEC_VPX,
            label: "VPX"
        },
    ]

    return <TerminalMetadataContainer>
        {/* Terminal Name */}
        <ConfigItemTitleContainer>
            <ConfigItemTitle>Terminal Name</ConfigItemTitle>
        </ConfigItemTitleContainer>
        <TerminalNameConfigContainer><TextField
            value={CurrentSelectedTerminal.name}
            disabled={CurrentSelectedTerminal.terminalConfigConfirm}
            onChange={handleTerminalNameUpdate}
            style={{zIndex: "0"}}
            id="terminal-name" 
            label="Terminal Name" 
            variant="outlined"
            size="small"
            fullWidth
        /></TerminalNameConfigContainer>

        {/* resolution */}
        <ConfigItemTitleContainer>
            <ConfigItemTitle>Resolution</ConfigItemTitle>
        </ConfigItemTitleContainer>
        <ResolutionConfigContainer>
            { /* 1280 x 720 (HD: 16:9) */}
            <ResolutionConfigEntry>
                <ResolutionShowcase 
                    Width={"256px"} 
                    Height={"144px"} 
                    Selected={CurrentSelectedTerminal.currentResolution === RESOLUTION_1280_720}
                    onClick={handleSelect_1280_720}
                >
                    1280 x 720
                </ResolutionShowcase>
            </ResolutionConfigEntry>

            { /* 640 x 480 (VGA 4:3) */}
            <ResolutionConfigEntry>
                <ResolutionShowcase 
                    Width={"256px"} 
                    Height={"192px"} 
                    Selected={CurrentSelectedTerminal.currentResolution === RESOLUTION_640_480}
                    onClick={handleSelect_640_480}
                >
                    640 x 480
                </ResolutionShowcase>
            </ResolutionConfigEntry>

            { /* 1024 x 600 (WIDTH 16:10) */}
            <ResolutionConfigEntry>
                <ResolutionShowcase 
                    Width={"204px"} 
                    Height={"120px"} 
                    Selected={CurrentSelectedTerminal.currentResolution === RESOLUTION_1024_600}
                    onClick={handleSelect_1024_600}
                >
                    1024 x 600
                </ResolutionShowcase>
            </ResolutionConfigEntry>

            { /* 1280 x 1024 (WIDTH 5:4) */}
            <ResolutionConfigEntry>
                <ResolutionShowcase 
                    Width={"256px"} 
                    Height={"204px"} 
                    Selected={CurrentSelectedTerminal.currentResolution === RESOLUTION_1280_1024}
                    onClick={handleSelect_1280_1024}
                >
                    1280 x 1024
                </ResolutionShowcase>
            </ResolutionConfigEntry>
        </ResolutionConfigContainer>

        {/* Frame Per Second (FPS) */}
        <ConfigItemTitleContainer>
            <ConfigItemTitle>Frame Per Second (FPS)</ConfigItemTitle>
        </ConfigItemTitleContainer>
        <TextField
          id="outlined-select-currency"
          select
          label="FPS"
          disabled={CurrentSelectedTerminal.terminalConfigConfirm}
          value={CurrentSelectedTerminal.currentFPS}
          onChange={handleChangeFPS}
          fullWidth
        >
          {fpsFormat.map((option) => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </TextField>

        {/* VCODEC */}
        <ConfigItemTitleContainer>
            <ConfigItemTitle>Video Codec</ConfigItemTitle>
        </ConfigItemTitleContainer>
        <TextField
          id="outlined-select-currency"
          select
          label="Video Codec"
          disabled={CurrentSelectedTerminal.terminalConfigConfirm}
          value={CurrentSelectedTerminal.vCodec}
          onChange={handleChangeVCodec}
          fullWidth
        >
          {vcodecFormat.map((option) => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </TextField>

        <ButtonGroupContainer>
            {
                !CurrentSelectedTerminal.terminalConfigConfirm && <Button 
                    variant="contained"
                    onClick={handleConfirmTerminalConfig}
                    sx={{width: "100%"}}
                >
                    Confirm Terminal Configuration
                </Button>
            }

            {
                CurrentSelectedTerminal.terminalConfigConfirm && <Button 
                    variant="contained"
                    onClick={handleCancelTerminalConfig}
                    color="success"
                    sx={{width: "100%"}}
                >
                    Cancel Terminal Configuration
                </Button>
            }
        </ButtonGroupContainer>
    </TerminalMetadataContainer> 


}

export default TerminalMetadata
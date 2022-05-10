import styled from 'styled-components';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../../Data/Reducers/terminalReducer';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import Tager from '../../../../Components/Tager/tager';
import { TabIndex_Dashboard_Terminal } from '../../terminalConfig';

const ApplicationDetailsContainer = styled.div`
    width: 100%;
    min-height: 300px;
    display: flex;
    margin: 10px 0px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
`

const ApplicationDetailsShowcase = styled.div`
    width: 95%;
    padding: 20px 20px;
    min-height: 300px;
    display: flex;
    flex-direction: column;
`

const ApplicationNameContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
`

const ApplicationMetadataContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: left;
`

const ApplicationDescriptionContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin-top: 30px;
`

const ButtonGroupContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 20px 0px;
`

export const APPLICATION_DETAIL_STATE_LOADED = "application_detail_state_loaded"
export const APPLICATION_DETAIL_STATE_LOADING = "application_detail_state_loading"

const ApplicationDetails = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    let StateTerminals = useSelector(state => state.terminal.StateTerminals)
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    let CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta

    const handleSelectApplication = (event) => {
        // confirm selected application
        dispatch(TerminalActions.updateTerminal({
            "type": "CONFIRM_SELECTED_APPLICATION",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))

        // jump to application terminal configuration panel
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_DASHBOARD_ID",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "dashboard_id": TabIndex_Dashboard_Terminal
        }))
    }

    const handleCancelSelectedApplication = (event) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "UNCONFIRM_SELECTED_APPLICATION",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))
    }

    return <ApplicationDetailsContainer>
        {/* load skeleton if application is not loaded  */}
        {
            CurrentApplicationMeta.currentSelectedApplicationDetailsState === APPLICATION_DETAIL_STATE_LOADING && <Skeleton 
                variant="rectangular"
                width={"100%"} 
                height={400}
            />
        }

        {/* load application info if application is loaded  */}
        {CurrentApplicationMeta.currentSelectedApplicationDetailsState ===  APPLICATION_DETAIL_STATE_LOADED && <ApplicationDetailsShowcase>
            <ApplicationNameContainer>
                {/* application name */}
                <Typography 
                    variant="h4" 
                    component="div" 
                    gutterBottom
                    sx={{fontFamily: "sans-serif"}}
                >
                    {CurrentApplicationMeta.currentSelectedApplication.name}
                </Typography>

                {/* application matedata */}
                <ApplicationMetadataContainer>
                    {/* os type */}
                    <Tager TagerConfig={{key: "Operating Systems", value: CurrentApplicationMeta.currentSelectedApplication.operatingSystem}} />

                    {/* usage count */}
                    <Tager TagerConfig={{key: "Usage Count", value: CurrentApplicationMeta.currentSelectedApplication.usageCount}} />
                </ApplicationMetadataContainer>
            </ApplicationNameContainer>

            {/* application description */}
            <ApplicationDescriptionContainer>
                <TextField
                    label="Application Description"
                    id="application-description"
                    value={CurrentApplicationMeta.currentSelectedApplication.description}
                    multiline
                    fullWidth
                    variant="outlined"
                    sx={{color: "#000000"}}
                />
            </ApplicationDescriptionContainer>

            {/* button group */}
            <ButtonGroupContainer>
                {
                    CurrentApplicationMeta.selectedApplicationConfirmed === false && <Button 
                            variant="contained"
                            onClick={handleSelectApplication}
                            sx={{width: "100%"}}
                        >
                            Select This Application
                        </Button>
                }

                {
                    CurrentApplicationMeta.selectedApplicationConfirmed === true && <Button 
                            variant="contained"
                            onClick={handleCancelSelectedApplication}
                            disabled={CurrentSelectedTerminal.wsConnectionStarted}
                            color="success"
                            sx={{width: "100%"}}
                        >
                            Cancel Selection
                        </Button>
                }
            </ButtonGroupContainer>
        </ApplicationDetailsShowcase>}

    </ApplicationDetailsContainer>
}

export default ApplicationDetails
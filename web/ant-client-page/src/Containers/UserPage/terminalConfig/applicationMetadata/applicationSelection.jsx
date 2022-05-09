import styled from 'styled-components';
import axios from 'axios'
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../../Data/Reducers/terminalReducer';
import TerminalImage from '../../../../Statics/Images/terminal.jpeg'
import DesktopImage from '../../../../Statics/Images/desktop.jpeg'
import { TABLE_STATE_LOADED, TABLE_STATE_LOADING } from '../../../../Components/RealtimeTable/realTimeTable';

const ApplicationSelectionContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
    margin: 20px 0px;
`

const BackgroundContainer = styled.div`
    background-image: url(${ ({Img}) => Img ? Img : null });
    background-size: cover;
    opacity: 0.7;
    position: absolute;
    border-radius: 20px;
    /* top: 0;
    left: 0; */
    width: 100%;
    height: 100%;
`

const ApplicationTypeContainer = styled.div`
    width: 95%;
    min-height: 180px;
    color: #ffffff;
    position: relative;
    border-radius: 20px;
    transition: 300ms;
    margin-bottom: 20px;

    &:hover {
        box-shadow: 5px 5px 5px #a0a0a0;
        cursor: pointer;
    }
`

const ApplicationTypeTitleContainer = styled.div`
    width: 80%;
    padding: 0px 20px;
    position: relative;
`

const ApplicationTypeTitle = styled.h1`
    font-size: 30px;
    margin-bottom: 0px;
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
    font-style: oblique;
`

const ApplicationTypeBriefContainer = styled.div`
    width: 50%;
    padding: 0px 20px;
    position: relative;
    background-color: #85858552;
    padding: 5px 10px;
    margin: 10px 0px;
    transition: 500ms;

    &:hover {
        background-color: #858585c7;
    }
`

const ApplicationTypeBrief = styled.p`
    font-size: 10px;
    color: #ffffff;
    margin-top: 5px;
    font-family: Verdana, Geneva, Tahoma, sans-serif;
    font-style: oblique;
    font-weight: bolder;
`

const ApplicationSelection = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    const StateTerminals = useSelector(state => state.terminal.StateTerminals)
    const StateAPI = useSelector(state => state.api.StateAPI)
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    let CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta

    /*
        @callback: HandleClickDesktop
        @description: handle click on desktop application
    */
    const HandleClickDesktop = () => {
        // update selected application
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_APPLICATION_TYPE",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "application_type": "stream"
        }))

        // check whether nav "app_market" exist, add if not
        let isExisted = false
        for(let i=0; i<CurrentApplicationMeta.applicationNavs.length; i++){
            if(CurrentApplicationMeta.applicationNavs[i].id === "app_market"){
                isExisted = true
                break
            }
        }
        if(!isExisted){
            dispatch(TerminalActions.updateTerminal({
                "type": "ADD_APPLICATION_NAV",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "new_nav": {
                    name: "App Market",
                    id: "app_market",
                }
            }))
        }

        // check whether nav "app_info" exist, delete if it's
        isExisted = false
        for(let i=0; i<CurrentApplicationMeta.applicationNavs.length; i++){
            if(CurrentApplicationMeta.applicationNavs[i].id === "app_info"){
                isExisted = true
                break
            }
        }
        if(isExisted){
            dispatch(TerminalActions.updateTerminal({
                "type": "DELETE_APPLICATION_NAV",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "nav_id": "app_info",
            }))
        }

        // update application list pagination
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_APPLICATION_PAGINATION",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "application_page": 1
        }))

        // clear application list
        dispatch(TerminalActions.updateTerminal({
            "type": "CLEAR_APPLICATION_LIST",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))

        // update state of application list
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_APPLICATION_LIST_STATE",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "application_list_state": TABLE_STATE_LOADING,
        }))

        // fetching total amount of applications
        axios.get(`${StateAPI.ScheduleProtocol}://${StateAPI.ScheduleHostAddr}:${StateAPI.SchedulePort}${StateAPI.ScheduleBaseURL}/${StateAPI.ScheduleAPI.ApplicationAmount}?type=${"stream"}`)
        .then((response) => {
            // update pagination count
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_APPLICATION_AMOUNT",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "application_amount": response.data.count,
            }))

            // fetching applications list (asynchronous)
            axios.get(`${StateAPI.ScheduleProtocol}://${StateAPI.ScheduleHostAddr}:${StateAPI.SchedulePort}${StateAPI.ScheduleBaseURL}/${StateAPI.ScheduleAPI.ApplicationList}?type=${"stream"}&page=${CurrentApplicationMeta.currentSelectedApplicationPageIndex}&size=${CurrentApplicationMeta.maxApplicationPerPage}&order=${CurrentApplicationMeta.applicationListOrderBy}`)
            .then((response) => {
                let fetchedApplicationList = response.data.applications

                // update application list for current terminal
                for(let i=0; i<fetchedApplicationList.length; i++){
                    let newListEntry = {
                        values: {
                            application_name: fetchedApplicationList[i].application_name,
                            application_id: fetchedApplicationList[i].application_id,
                            create_user: fetchedApplicationList[i].create_user,
                            updated_at: fetchedApplicationList[i].updated_at,
                            usage_count: fetchedApplicationList[i].usage_count
                        },
                        
                        selected: false
                    }

                    dispatch(TerminalActions.updateTerminal({
                        "type": "ADD_APPLICATION_LIST",
                        "terminal_key": `${StateTerminals.currentSelected}`,
                        "application": newListEntry,
                    }))
                }

                // update state of application list
                dispatch(TerminalActions.updateTerminal({
                    "type": "UPDATE_APPLICATION_LIST_STATE",
                    "terminal_key": `${StateTerminals.currentSelected}`,
                    "application_list_state": TABLE_STATE_LOADED,
                }))
            })
            .catch((error) => {
                // found error
                if(error.response){
                    console.log(error.response)
                    dispatch(TerminalActions.updateTerminal({
                        "type": "UPDATE_APPLICATION_LIST_STATE",
                        "terminal_key": `${StateTerminals.currentSelected}`,
                        "application_list_state": TABLE_STATE_LOADED,
                    }))
                }
            })
        })
        .catch((error) => {
            if(error.response){
                console.log(error.response)
            }
        })

        // update current nav
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_SELECTED_NAV",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "current_selected_nav": "app_market"
        }))
    }

    return (
        <ApplicationSelectionContainer>
            {/* Desktop Application */}
            <ApplicationTypeContainer
                onClick={HandleClickDesktop}
            >
                <BackgroundContainer Img={DesktopImage} />
                <ApplicationTypeTitleContainer><ApplicationTypeTitle>Desktop Application</ApplicationTypeTitle></ApplicationTypeTitleContainer>
                <ApplicationTypeBriefContainer><ApplicationTypeBrief>
                    Desktop applictaions provide human-friendly interation interface, commonly used for cloud gaming, engineering designing, etc.
                </ApplicationTypeBrief></ApplicationTypeBriefContainer>
            </ApplicationTypeContainer>

            {/* Console Application */}
            <ApplicationTypeContainer>
                <BackgroundContainer Img={TerminalImage} />
                <ApplicationTypeTitleContainer><ApplicationTypeTitle>Console Application</ApplicationTypeTitle></ApplicationTypeTitleContainer>
                <ApplicationTypeBriefContainer><ApplicationTypeBrief>
                    Console applictaions provide terminal-style interaction with remote clients, commonly used by software developers who have no need of desktop user interface for operating.
                </ApplicationTypeBrief></ApplicationTypeBriefContainer>
            </ApplicationTypeContainer>
        </ApplicationSelectionContainer>
    )
}

export default ApplicationSelection
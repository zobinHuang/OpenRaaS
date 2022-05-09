import styled from 'styled-components';
import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../Data/Reducers/terminalReducer';
import AppsIcon from '@mui/icons-material/Apps';
import SwipeLeftIcon from '@mui/icons-material/SwipeLeft';
import InfoIcon from '@mui/icons-material/Info';
import BreadCrumbsNav from '../../../Components/BreadCurmbsNav/breadCurmbsNav';
import ApplicationSelection from './applicationMetadata/applicationSelection';
import ApplicationList from './applicationMetadata/applicationList';
import ApplicationDetails from './applicationMetadata/applicationDetails';

const ApplicationMetadataContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 10px 20px;
    overflow-y: hidden;
`

const ApplicationMetadata = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    let StateTerminals = useSelector(state => state.terminal.StateTerminals)
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    let CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta
    
    /*
        @callback: HandleClickNavButton
        @description: Handle click on nav button
    */
    const HandleClickNavButton = (navID) => {
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_SELECTED_NAV",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "current_selected_nav": navID
        }))
    }

    /*
        @config: BreadCrumbsNavConfig
        @description: Config props for component BreadCrumbsNav
    */
    const BreadCrumbsNavConfig = {
        currentSelectedNavID: CurrentApplicationMeta.currentSelectedApplicationNav,
        handleClick: HandleClickNavButton,
        navIcons: {
            "app_type": <SwipeLeftIcon />,
            "app_market": <AppsIcon />,
            "app_info": <InfoIcon />
        },
        navs: CurrentApplicationMeta.applicationNavs,
        position: "left",
        disabled: CurrentApplicationMeta.selectedApplicationConfirmed,
    }

    return (
        <ApplicationMetadataContainer>
            {/* BreadCrumbs Navigation */}
            <BreadCrumbsNav BreadCrumbsNavConfig={BreadCrumbsNavConfig} />

            {/* Subpage: Application Type Selection */}    
            {CurrentApplicationMeta.currentSelectedApplicationNav === "app_type" && <ApplicationSelection />}

            {/* Subpage: Realtime Table */}
            {CurrentApplicationMeta.currentSelectedApplicationNav === "app_market" && <ApplicationList />}

            {/* Subpage: Application Details */}
            {CurrentApplicationMeta.currentSelectedApplicationNav === "app_info" && <ApplicationDetails />}
        </ApplicationMetadataContainer>
    )
}

export default ApplicationMetadata
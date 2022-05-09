import { useDispatch, useSelector } from 'react-redux';
import { actions as TerminalActions } from '../../../../Data/Reducers/terminalReducer';
import RealTimeTable, { TABLE_STATE_LOADED } from '../../../../Components/RealtimeTable/realTimeTable';
import axios from 'axios'
import { APPLICATION_DETAIL_STATE_LOADED } from './applicationDetails';
export const APP_LIST_ORDER_BY_NAME = "orderByName"
export const APP_LIST_ORDER_BY_UPDATE_TIME = "orderByUpdateTime"
export const APP_LIST_ORDER_BY_USAGE_COUNT = "orderByUsageCount"

const ApplicationList = (props) => {
    // get dispatch
    const dispatch = useDispatch()

    // get global states of terminal reducer
    let StateTerminals = useSelector(state => state.terminal.StateTerminals)
    let CurrentSelectedTerminal = StateTerminals.terminalsMap[StateTerminals.currentSelected]
    let CurrentApplicationMeta = StateTerminals.terminalsMap[StateTerminals.currentSelected].applicationMeta
    let StateAPI = useSelector(state => state.api.StateAPI)

    /*
        @callback: handleChangePagination
        @description: handle change pagination
    */
    const handleChangePagination= (event, page) => {
        // clear application list
        dispatch(TerminalActions.updateTerminal({
            "type": "CLEAR_APPLICATION_LIST",
            "terminal_key": `${StateTerminals.currentSelected}`,
        }))

        // change pagination selection
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_APPLICATION_PAGINATION",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "application_page": page
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
            axios.get(`${StateAPI.ScheduleProtocol}://${StateAPI.ScheduleHostAddr}:${StateAPI.SchedulePort}${StateAPI.ScheduleBaseURL}/${StateAPI.ScheduleAPI.ApplicationList}?type=${"stream"}&page=${page}&size=${CurrentApplicationMeta.maxApplicationPerPage}&order=${CurrentApplicationMeta.applicationListOrderBy}`)
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
    }

    /*
        @callback: handleClickOnRow
        @description: handle click on table row
    */
    const handleClickOnRow = (event) => {
        // extract selected row
        let selectedRow = event.target.id

        // record selected application id to terminal store
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_SELECTED_APPLICATION",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "list_index": selectedRow
        }))

        // check whether nav "app_info" exist, add if it isn't
        let isExisted = false
        for(let i=0; i<CurrentApplicationMeta.applicationNavs.length; i++){
            if(CurrentApplicationMeta.applicationNavs[i].id === "app_info"){
                isExisted = true
                break
            }
        }
        if(!isExisted){
            dispatch(TerminalActions.updateTerminal({
                "type": "ADD_APPLICATION_NAV",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "new_nav": {
                    name: "App Details",
                    id: "app_info",
                }
            }))
        }

        // request selected application details
        axios.get(`${StateAPI.ScheduleProtocol}://${StateAPI.ScheduleHostAddr}:${StateAPI.SchedulePort}${StateAPI.ScheduleBaseURL}/${StateAPI.ScheduleAPI.ApplicationDetails}?id=${CurrentApplicationMeta.applicationList[selectedRow].values.application_id}&type=${CurrentApplicationMeta.currentSelectedApplicationType}`)
        .then((response) => {
            let fetchedApplication = response.data.application

            // update application
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_APPLICATION_DETAILS",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "application_id": CurrentApplicationMeta.applicationList[selectedRow].values.application_id,
                "application_name": fetchedApplication.application_name,
                "application_creator": fetchedApplication.create_user,
                "application_create_time": fetchedApplication.create_at,
                "application_update_time": fetchedApplication.updated_at,
                "application_description": fetchedApplication.description,
                "application_operating_system": fetchedApplication.operating_system,
                "application_usage_count": fetchedApplication.usage_count,
            }))

            // update application detail state
            dispatch(TerminalActions.updateTerminal({
                "type": "UPDATE_APPCATION_DETAILS_STATE",
                "terminal_key": `${StateTerminals.currentSelected}`,
                "application_details_state": APPLICATION_DETAIL_STATE_LOADED,
            }))
        })
        .catch((error) => {
            // found error
            if(error.response){
                console.log(error.response)
                
                // TODO: do more prompt
            }
        })

        // update current nav
        dispatch(TerminalActions.updateTerminal({
            "type": "UPDATE_SELECTED_NAV",
            "terminal_key": `${StateTerminals.currentSelected}`,
            "current_selected_nav": "app_info"
        }))
    }
    
    /*
        @callback: HandleOrderByApplicationName
        @description: Handle click on order by application name
    */
    const HandleOrderByApplicationName = (event) => {
        // TODO
        console.log("HandleOrderByApplicationName")
    }

    /*
        @callback: HandleOrderByUpdateTime
        @description: Handle click on order by update time
    */
    const HandleOrderByUpdateTime = (event) => {
        // TODO
        console.log("HandleOrderByUpdateTime")
    }

    /*
        @callback: HandleOrderByUsageCount
        @description: Handle click on order by usage count
    */
    const HandleOrderByUsageCount = (event) => {
        // TODO
        console.log("HandleOrderByUsageCount")
    }

    /*
        @config: RealTimeTableConfig
        @description: Config props for component RealTimeTable
    */
    const RealTimeTableConfig = {
        // -------------------- global related ------------------
        // disable RealTimeTable component
        "disabled": CurrentSelectedTerminal.currentStepIndex !== 0,

        // -------------------- search bar related ------------------
        // seach bar place holder
        "searchBarPlaceHolder": "Search Applications",

        // -------------------- table related ------------------
        // table head
        "heads": [
            {
                "index": "application",
                "name": "Application",
                "canOrder": true,
                "OrderCallback": HandleOrderByApplicationName
            },
            {
                "index": "application_id",
                "name": "Application ID",
            },
            {
                "index": "creator",
                "name": "Creator",
            },
            {
                "index": "last_update",
                "name": "Last Update",
                "canOrder": true,
                "OrderCallback": HandleOrderByUpdateTime
            },
            {
                "index": "usage_count",
                "name": "Usage Count",
                "canOrder": true,
                "OrderCallback": HandleOrderByUsageCount
            }
        ],

        // table rows
        "rows": CurrentApplicationMeta.applicationList,

        // callback: handle double click on row
        "handleClickOnRow": handleClickOnRow,

        // overall potential rows amount
        "rowAmount": CurrentApplicationMeta.applicationAmount,       

        // maximum rows displayed on one page
        "maxRowsPerPage": CurrentApplicationMeta.maxApplicationPerPage,   

        //
        // application list state
        "tableState": CurrentApplicationMeta.applicationListState,

        // -------------------- pagination related ------------------
        // current page selection
        "currentSelectedPageIndex": CurrentApplicationMeta.currentSelectedApplicationPageIndex,

        // handle change pagination
        "handleChangePagination": handleChangePagination,
    }

    return <RealTimeTable RealTimeTableConfig={RealTimeTableConfig} />
}

export default ApplicationList
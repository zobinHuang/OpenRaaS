import { createSlice } from '@reduxjs/toolkit'
import { TABLE_STATE_LOADED, TABLE_STATE_LOADING } from '../../Components/RealtimeTable/realTimeTable'
import { TabIndex_Dashboard_LogViewer } from '../../Containers/UserPage/terminalConfig'
import APIConfig from "../../Configurations/APIConfig.json"
import axios from 'axios'
import GetUUID from '../../Utils/get_uuid'
import { APP_LIST_ORDER_BY_NAME, APP_LIST_ORDER_BY_UPDATE_TIME, APP_LIST_ORDER_BY_USAGE_COUNT } from '../../Containers/UserPage/terminalConfig/applicationMetadata/applicationList'
import { APPLICATION_DETAIL_STATE_LOADED, APPLICATION_DETAIL_STATE_LOADING } from '../../Containers/UserPage/terminalConfig/applicationMetadata/applicationDetails'
import { FPS_30, RESOLUTION_1280_720 } from '../../Containers/UserPage/terminalConfig/terminalMetadata'

const defaultTerminalState = {
    // Name of the terminal
    name: "New Terminal", 

    // Steps of a terminal
    steps: [
        {
            name: "Config Instance",
            id: "config",
            state: "inStep",    // beforeStep, inStep, afterStep, failedStep
            descriptionMessage: "Waiting for finishing local configuration",
            failedMessage: "Incorrect Configuration Parameters"
        },
        {
            name: "Connect to Scheduler",
            id: "connectScheduler",
            state: "beforeStep",
            descriptionMessage: "Connecting to scheduler",
            failedMessage: "Failed to Connect to Scheduler"
        },
        {
            name: "Schedule Compute Node",
            id: "scheduleComp",
            state: "beforeStep",
            descriptionMessage: "The schduler is searching for compute node",
            failedMessage: "Failed to Schedule Compute Node"
        },
        {
            name: "Schedule Storage Node",
            id: "scheduleStorage",
            state: "beforeStep",
            descriptionMessage: "The provider is searching for storage node",
            failedMessage: "Failed to Schedule Storage Node"
        },
        {
            name: "Prepare Instance",
            id: "prepare",
            state: "beforeStep",
            descriptionMessage: "The compute node is preparing instance",
            failedMessage: "Failed to Prepare Instance on Compute Node"
        },
        {
            name: "Run Instance",
            id: "run",
            state: "beforeStep",
            descriptionMessage: "The instance is now ready for operating",
            failedMessage: "Failed to Run Instance on Compute Node"
        }
    ],

    // Log infomation of current terminal
    logInfo: [
        // {
        //     "priority": "ERROR", // (ERROR, INFO, WARN)
        //     "time": ""
        //     "content": "error message",
        // },
    ],

    // instance id in scheduler
    instanceSchedulerID: "",

    // amount of newly unread log
    unreadLogCount: 0,

    // level of unread log
    unreadLogLevel: "primary", // primary / error

    // current dashboard index
    currentDashboardIndex: 0,

    // current step of the terminal (0: config, 1: connectScheduler, 2: scheduleComp, 3: scheduleStorage, 4: run)  
    currentStepIndex: 0,
    
    // ice server list
    iceServers: [],

    // unique client ID in scheduler
    clientID: "",

    // state to inidicate whether websocket connection has started
    wsConnectionStarted: false,

    // current selected resolution
    currentResolution: RESOLUTION_1280_720,
    screenHeight: 0,
    screenWidth: 0,

    // current frame per second (fps)
    currentFPS: FPS_30,

    // confirmation of terminal metadata has been configured
    terminalConfigConfirm: false,

    // application metadata for current terminal
    applicationMeta: {
        // application nav for current terminal
        applicationNavs: [
            {
                name: "App Type",
                id: "app_type",
            }
        ],

        // current selected application type
        currentSelectedApplicationType: "",  // stream / console

        // index of current selected nav
        currentSelectedApplicationNav: "app_type",

        // application search result for current terminal
        applicationList: [
            // {
            //     values: {
            //         application_name: "road rash",
            //         application_id: "xxxx-xxxx-xxxx-xxxx",
            //         create_user: "zobinHuang",
            //         updated_at: "Feb.14 2022",
            //         usage_count: 2013,
            //     },
            //     selected: false,
            // }
        ],

        // application list state for current terminal
        applicationListState: TABLE_STATE_LOADED,

        // current selected application list index
        currentSelectedApplicationListIndex: -1,

        // application pagination for current terminal
        currentSelectedApplicationPageIndex: 1,

        // total application amount, used for calculating pagination count for current terminal
        applicationAmount: 1, 

        // order scheme of application list for current terminal
        applicationListOrderBy: APP_LIST_ORDER_BY_NAME,

        // maximum showed application per page
        maxApplicationPerPage: 10,

        // application details state for current terminal
        currentSelectedApplicationDetailsState: APPLICATION_DETAIL_STATE_LOADING,

        // application details
        currentSelectedApplication: {
            // application meta
            name: "",
            id: "",
            creator: "",
            createTime: "",
            updateTime: "",
            description: "",

            // enviroment meta
            operatingSystem: "",

            // platform meta
            usageCount: ""
        },

        // confirmation of application has been selected
        selectedApplicationConfirmed: false
    }
}

const terminalsSlice = createSlice({
    name: 'terminals',
    
    initialState: {
        StateTerminals : {
            terminalsMap: {},
            currentSelected: "",
            currentSelectedIndex: 0
        }
    },
    
    reducers: {
        /* Add a new terminal to the list */
        addTerminal(state,action){
            // record current map length as new terminal index
            // const map_length = Object.getOwnPropertyNames(state.StateTerminals.terminalsMap).length

            // generate terminal uuid
            let terminalUUID = GetUUID()

            // deep clone object
            let newTerminalState = structuredClone(defaultTerminalState)
            newTerminalState.name = `New Terminal`

            // insert into list
            let newTerminalKey = terminalUUID
            state.StateTerminals.terminalsMap[newTerminalKey] = newTerminalState

            // change current selected index
            state.StateTerminals.currentSelected = newTerminalKey
            state.StateTerminals.currentSelectedIndex = Object.getOwnPropertyNames(state.StateTerminals.terminalsMap).length-1
        },

        /* Delete the specified terminal from the list */
        deleteTerminal(state,action){
            // delete selected terminal from terminal map
            delete state.StateTerminals.terminalsMap[action.payload]

            // change current selected terminal
            if(state.StateTerminals.currentSelectedIndex == 0){
                state.StateTerminals.currentSelectedIndex = 0
                state.StateTerminals.currentSelected = Object.keys(state.StateTerminals.terminalsMap)[0]
            }
            else {
                state.StateTerminals.currentSelectedIndex -= 1
                state.StateTerminals.currentSelected = Object.keys(state.StateTerminals.terminalsMap)[state.StateTerminals.currentSelectedIndex]
            }
        },

        /* Change the information of specified terminal */
        updateTerminal(state,action){
            if(!state.StateTerminals.terminalsMap[action.payload.terminal_key]){
                console.log("Warn - unknown terminal while updateTerminal, abandoned")
                return
            }
            switch(action.payload.type){                
                /* Case: update terminal name */
                case "UPDATE_TERMINAL_NAME":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].name = action.payload.name
                    break;

                /* Case: update ice server list */
                case "UPDATE_ICESERVERS":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].iceServers = action.payload.ice_servers
                    break;

                /* Case: update dashboard index */
                case "UPDATE_DASHBOARD_ID":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].currentDashboardIndex = action.payload.dashboard_id
                    break;

                /* Case: update client index */
                case "UPDATE_CLIENT_ID":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].clientID = action.payload.client_id
                    break;

                /* Case: update unread log count */
                case "UPDATE_UNREAD_LOG_COUNT":
                    // update unread log count
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].unreadLogCount = action.payload.unread_log_count

                    // reset unread log level
                    if(action.payload.unread_log_count === 0){
                        state.StateTerminals.terminalsMap[action.payload.terminal_key].unreadLogLevel = "primary"
                    }
                    break;

                /* Case: update terminal current step */
                case "UPDATE_TERMINAL_STEP":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].currentStepIndex = action.payload.current_step_index
                    break;
                
                /* Case: update terminal resolution */
                case "UPDATE_INSTANCE_SCHEDULER_ID":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].instanceSchedulerID = action.payload.instance_scheduler_id
                    break;
                
                /* Case: update terminal resolution */
                case "UPDATE_TERMINAL_RESOLUTION":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].currentResolution = action.payload.resolution
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].screenHeight = action.payload.height
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].screenWidth = action.payload.width
                    break;
                
                /* Case: update terminal resolution */
                case "UPDATE_TERMINAL_FPS":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].currentFPS = action.payload.fps
                    break;

                /* Case: update selected application type */
                case "UPDATE_APPLICATION_TYPE":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplicationType = action.payload.application_type
                    break;
                
                /* Case: update total application amount */
                case "UPDATE_APPLICATION_AMOUNT":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationAmount = action.payload.application_amount
                    break;
                
                /* Case: confirm websocket connection started */
                case "CONFIRM_WS_CONNECTION_STARTED":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].wsConnectionStarted = true
                    break;

                /* Case: unconfirm websocket connection started */
                case "UNCONFIRM_WS_CONNECTION_STARTED":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].wsConnectionStarted = false
                    break;

                /* Case: confirm selected application */
                case "CONFIRM_SELECTED_APPLICATION":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.selectedApplicationConfirmed = true
                    break;

                /* Case: unconfirm selected application */
                case "UNCONFIRM_SELECTED_APPLICATION":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.selectedApplicationConfirmed = false
                    break;

                /* Case: confirm terminal configuration */
                case "CONFIRM_TERMINAL_CONFIG":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].terminalConfigConfirm = true
                    break;

                /* Case: unconfirm terminal configuration */
                case "UNCONFIRM_TERMINAL_CONFIG":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].terminalConfigConfirm = false
                    break;

                /* Case: update state of selected application details */
                case "UPDATE_APPCATION_DETAILS_STATE":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplicationDetailsState = action.payload.application_details_state
                    break;

                /* Case: update application selected nav */
                case "UPDATE_SELECTED_NAV":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplicationNav = action.payload.current_selected_nav
                    break;

                /* Case: add application nav */
                case "ADD_APPLICATION_NAV":
                    let add_nav_length = state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationNavs.length
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationNavs[add_nav_length] = action.payload.new_nav
                    break;
                
                /* Case: delete application nav */
                case "DELETE_APPLICATION_NAV":
                    for(let i=0; i<state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationNavs.length; i++){
                        if(state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationNavs[i].id === action.payload.nav_id){
                            state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationNavs.splice(i,1)
                            break
                        } 
                    }
                    break;

                /* Case: update selected update application pagination */
                case "UPDATE_APPLICATION_PAGINATION":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplicationPageIndex = action.payload.application_page
                    break;
                
                /* Case: update selected application */
                case "UPDATE_SELECTED_APPLICATION":
                    // update outsided record selected index
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplicationListIndex = action.payload.list_index
                    
                    // clear prev state
                    for(let i=0; i<state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList.length; i++){
                        if(state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList[i].selected === true){
                            state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList[i].selected = false
                        }
                    }

                    // set new state
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList[action.payload.list_index].selected = true
                    break;
                
                /* Case: add application entry into application list */
                case "CLEAR_APPLICATION_LIST":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList.splice(0)
                    break
                
                /* Case: add application entry into application list */
                case "ADD_APPLICATION_LIST":
                    let list_length = state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList.length
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationList[list_length] = action.payload.application
                    break

                /* Case: update state of application list */
                case "UPDATE_APPLICATION_LIST_STATE":
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.applicationListState = action.payload.application_list_state
                    break;
                
                /* Case: update details of current selected application */
                case "UPDATE_APPLICATION_DETAILS":
                    // update application name
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.id = action.payload.application_id

                    // update application name
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.name = action.payload.application_name

                    // update application creator
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.creator = action.payload.application_creator

                    // update application create time
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.createTime = action.payload.application_create_time

                    // update application update time
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.createTime = action.payload.application_update_time

                    // update application description
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.description = action.payload.application_description

                    // update application operating system
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.operatingSystem = action.payload.application_operating_system

                    // update application usage count
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].applicationMeta.currentSelectedApplication.usageCount = action.payload.application_usage_count
                    break;

                /* Case: append log */
                case "APPEND_LOG_CONTENT":
                    // append log
                    let log_priority = action.payload.log_priority
                    let log_time = action.payload.log_time
                    let log_content = action.payload.log_content
                    let log_list_length = state.StateTerminals.terminalsMap[action.payload.terminal_key].logInfo.length
                    state.StateTerminals.terminalsMap[action.payload.terminal_key].logInfo[log_list_length] = ({
                        "priority": log_priority,
                        "time": log_time,
                        "content": log_content,
                    })

                    // update unread log count
                    if(state.StateTerminals.terminalsMap[action.payload.terminal_key].currentDashboardIndex !== TabIndex_Dashboard_LogViewer){
                        state.StateTerminals.terminalsMap[action.payload.terminal_key].unreadLogCount += 1
                    }

                    // set unread log level if necessary
                    if(log_priority === "ERROR"){
                        state.StateTerminals.terminalsMap[action.payload.terminal_key].unreadLogLevel = "error"
                    }

                    break;
            }
        },

        /* Change the selected terminal tab */
        changeSelectedTab(state,action){
            state.StateTerminals.currentSelectedIndex = action.payload
            state.StateTerminals.currentSelected = Object.keys(state.StateTerminals.terminalsMap)[state.StateTerminals.currentSelectedIndex]
        },

        /* Reset state */
        reset(state,action){
            state.StateTerminals.terminalsList = []
            state.StateTerminals.currentSelectedIndex = 0
        }
    }
})

export const actions = terminalsSlice.actions
export default terminalsSlice.reducer

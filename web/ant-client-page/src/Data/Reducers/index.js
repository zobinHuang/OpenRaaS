import { combineReducers } from "redux";
import apiReducer from "./apiReducer";
import authReducer from "./authReducer";
import backdropReducer from "./backdropReducer";
import infoReducer from "./infoReducer";
import snackBarReducer from "./snackBarReducer";
import terminalReducer from "./terminalReducer"
import userSidebarReducer from "./userSidebarReducer"

export default combineReducers({
    // Info State
    info: infoReducer,

    // API State
    api: apiReducer,

    // Terminals State
    terminal: terminalReducer,

    // User Authentication State
    auth: authReducer,

    // Snackbar State
    snackbar: snackBarReducer,

    // Backdrop State
    backdrop: backdropReducer,

    // User Sidebar State
    userSidebar: userSidebarReducer
})
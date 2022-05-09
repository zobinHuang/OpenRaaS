import { combineReducers } from "redux";
import apiReducer from "./apiReducer";
import infoReducer from "./infoReducer";

export default combineReducers({
    // AuthAPI State
    api: apiReducer,

    // Info State
    info: infoReducer
})
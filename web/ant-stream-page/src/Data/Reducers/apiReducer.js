import { createSlice } from '@reduxjs/toolkit'
import APIConfig from "../../Configurations/APIConfig.json"

const apiSlice = createSlice({
    name: 'api',
    
    initialState: {
        StateAPI : APIConfig
    },
    
    reducers: {
        /* Leave Blank */
    }
})

export const actions = apiSlice.actions
export default apiSlice.reducer

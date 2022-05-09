import { createSlice } from '@reduxjs/toolkit'
import InfoConfig from "../../Configurations/InfoConfig.json"

const infoSlice = createSlice({
    name: 'info',
    
    initialState: {
        StateInfo : InfoConfig
    },
    
    reducers: {
        /* Leave Blank */
    }
})

export const actions = infoSlice.actions
export default infoSlice.reducer

import { createSlice } from '@reduxjs/toolkit'

const userSidebarSlice = createSlice({
    name: 'userSidebar',
    
    initialState: {
        StateUserSidebar: {
            showSidebar: true,
            currentSelected: "my_cloud"
        }
    },
    
    reducers: {
        /* Chenge the item that currently selected */
        changeCurrentSelected(state,action){
            state.StateUserSidebar.currentSelected = action.payload
        },
        
        /* show user sidebar */
        changeShow(state,action){
            if(action.payload == true || action.payload == false)
                state.StateUserSidebar.showSidebar = action.payload
            else
                state.StateUserSidebar.showSidebar = false
        },

        /* Reset state */
        reset(state,action){
            state.StateUserSidebar.showSidebar = true
            state.StateUserSidebar.currentSelected = "my_cloud"
        }
    }
})

export const actions = userSidebarSlice.actions
export default userSidebarSlice.reducer

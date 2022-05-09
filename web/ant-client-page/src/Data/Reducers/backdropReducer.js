import { createSlice } from '@reduxjs/toolkit'

const backdropSlice = createSlice({
    name: 'backdrop',
    
    initialState: {
        backdropEnabled: false
    },
    
    reducers: {
        /* Open backdrop */
        openBackdrop(state, action){
            state.backdropEnabled = true
        },

        /* Close backdrop */
        closeBackdrop(state, action){
            state.backdropEnabled = false
        }
    }
})

export const actions = backdropSlice.actions
export default backdropSlice.reducer

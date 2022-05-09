import { createSlice } from '@reduxjs/toolkit'

const snackBarSlice = createSlice({
    name: 'snackbar',
    
    initialState: {
        StateSnackBar: {
            snackBarEnabled: false,
            snackBarContent: ""
        }
    },
    
    reducers: {
        /* Open Snack Bar */
        showSnackBar(state,action){
            state.StateSnackBar.snackBarEnabled = true
            state.StateSnackBar.snackBarContent = action.payload
        },

        /* Close Snack Bar */
        closeSnackBar(state,action){
            state.StateSnackBar.snackBarEnabled = false
        }
    }
})

export const actions = snackBarSlice.actions
export default snackBarSlice.reducer

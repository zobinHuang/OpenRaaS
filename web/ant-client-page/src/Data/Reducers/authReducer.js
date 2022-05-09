import { createSlice } from '@reduxjs/toolkit'

const authSlice = createSlice({
    name: 'auth',
    
    initialState: {
        StateAuth: {
            isLogin: false,
            username: "",
            idToken: "",
            refreshToken: "",
            avatarSrc: ""
        }
    },
    
    reducers: {
        /* Change login state */
        login(state,action){
            state.StateAuth.username = action.payload.username
            state.StateAuth.idToken = action.payload.idToken
            state.StateAuth.refreshToken = action.payload.refreshToken
            state.StateAuth.avatarSrc = action.payload.avatarSrc
            state.StateAuth.isLogin = true
        },
        
        /* Change login state */
        logout(state,action){
            state.StateAuth.username = ""
            state.StateAuth.idToken = ""
            state.StateAuth.refreshToken = ""
            state.StateAuth.avatarSrc = ""
            state.StateAuth.isLogin = false
        }
    }
})

export const actions = authSlice.actions
export default authSlice.reducer

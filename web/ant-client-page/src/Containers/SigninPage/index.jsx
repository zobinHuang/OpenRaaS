import React, { useState } from 'react';
import styled from 'styled-components';
import axios from 'axios'
import Slide from '@mui/material/Slide';
import FormFormat from "./form_format.json"
import { Navigate } from "react-router-dom";
import { useSelector, useDispatch } from 'react-redux'
import { actions as AuthActions } from '../../Data/Reducers/authReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import { actions as BackdropActions } from '../../Data/Reducers/backdropReducer'
import PageHeader from '../../Components/Header/header';
import PageFooter from '../../Components/Footer/footer';
import UserForm from '../../Components/UserForm/user_form';

const OuterContainer = styled.div`
    width: 100%;
    margin: 0px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    padding-top: calc(10vh);
    align-items: center;
    justify-content: center;
`

const SigninPageContainer = styled.div`
    width: 100%;
    height: calc(80vh);
    margin: 0px;
    padding: 0px;
    overflow-x: hidden;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const SigninPage = (props) => {
    const dispatch = useDispatch()

    // Get global state
    const StateAuth = useSelector(state => state.auth.StateAuth)
    const StateInfo = useSelector(state => state.info.StateInfo)
    const StateAPI = useSelector(state => state.api.StateAPI)

    // Define entry state
    const [email, changeEmail] = useState("")
    const [password, changePassword] = useState("")
    
    // Bind state to each entry
    FormFormat.entries[0].state = email
    FormFormat.entries[0].changeState = changeEmail
    FormFormat.entries[1].state = password
    FormFormat.entries[1].changeState = changePassword
    
    // Define callback
    const handleSignInClicked = async (event) => {
        event.preventDefault();
        
        // check empty
        if(email == ""){
            dispatch(SnackBarActions.showSnackBar("Please input your email"))
            return
        }
        if(password == ""){
            dispatch(SnackBarActions.showSnackBar("Please input your password"))
            return
        }

        /* Show Backdrop*/
        dispatch(BackdropActions.openBackdrop())
        // Config.BackdropConfig.ChangeBackdropEnabled(true)

        /* Login */
        axios.post(`${StateAPI.AuthProtocol}://${StateAPI.AuthHostAddr}:${StateAPI.AuthPort}${StateAPI.AuthBaseURL}/${StateAPI.AuthAPI.SignIn}`, {
            "email": email,
            "password": password
        })
        .then((response) => {
            /* Disabled Backdrop */
            dispatch(BackdropActions.closeBackdrop())
            // Config.BackdropConfig.ChangeBackdropEnabled(false)

            /* Change State */
            dispatch(AuthActions.login({
                username: response.data.username,
                idToken: response.data.tokens.idToken,
                refreshToken: response.data.tokens.refreshToken,
                avartarSrc: ""
            }))

            /* Show snack bar */
            dispatch(SnackBarActions.showSnackBar("Successfully Signin!"))
        })
        .catch((error) => {
            if(error.response){
                /* Disabled Backdrop */
                dispatch(BackdropActions.closeBackdrop())
                // Config.BackdropConfig.ChangeBackdropEnabled(false)

                /* Authentication error */
                if(error.response.status === 401){
                    dispatch(SnackBarActions.showSnackBar("Error Authentication! Check email and password."))
                }
            }else if(error.request){
                /* Disabled Backdrop */
                dispatch(BackdropActions.closeBackdrop())
                // Config.BackdropConfig.ChangeBackdropEnabled(false)
                
                /* No response */
                dispatch(SnackBarActions.showSnackBar("Can't connect to server, check your network!"))
            }
            return;
        })
    }

    /* Bind callback for button */
    FormFormat.buttons[0].callback = handleSignInClicked

    if(StateAuth.isLogin){
        return <Navigate to="/user" replace={true} />
    } else {
        return (
            <OuterContainer>
                <PageHeader />
                    <Slide 
                        direction={"right"} 
                        in={true}
                        mountOnEnter 
                        unmountOnExit
                    >
                        <SigninPageContainer>
                            <UserForm FormFormat={FormFormat}/>
                        </SigninPageContainer>
                    </Slide>
                <PageFooter Stick={true} />
            </OuterContainer>
        )
    }
}

export default SigninPage;

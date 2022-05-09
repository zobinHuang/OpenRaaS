import React, { useState } from 'react';
import axios from 'axios'
import styled from 'styled-components';
import { Navigate } from "react-router-dom";
import { useSelector, useDispatch } from 'react-redux'
import Slide from '@mui/material/Slide';
import { actions as AuthActions } from '../../Data/Reducers/authReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'
import { actions as BackdropActions } from '../../Data/Reducers/backdropReducer'
import FormFormat from './form_format.json'
import UserForm from '../../Components/UserForm/user_form';
import PageHeader from '../../Components/Header/header';
import PageFooter from '../../Components/Footer/footer';

const OuterContainer = styled.div`
    width: 100%;
    padding-top: calc(10vh);
    margin: 0px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const SignupPageContainer = styled.div`
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

const SignupPage = (props) => {    
    const dispatch = useDispatch()

    // Get global state
    const StateAuth = useSelector(state => state.auth.StateAuth)
    const StateInfo = useSelector(state => state.info.StateInfo)
    const StateAPI = useSelector(state => state.api.StateAPI)

    /* Private state to record whether signup success */
    const [signupState, changeSignupState] = useState(false) 

    /* Define entry state*/
    const [email, changeEmail] = useState("")
    const [password, changePassword] = useState("")
    const [confirmPassword, changeConfirmPassword] = useState("")
    
    /* Bind state to each entry*/
    FormFormat.entries[0].state = email
    FormFormat.entries[0].changeState = changeEmail
    FormFormat.entries[1].state = password
    FormFormat.entries[1].changeState = changePassword
    FormFormat.entries[2].state = confirmPassword
    FormFormat.entries[2].changeState = changeConfirmPassword

    /* Define callback */
    const handleSignUpClicked = async (event) => {
        event.preventDefault();

        /* check empty */
        if(email == ""){
            dispatch(SnackBarActions.showSnackBar("Please input your email"))
            return
        }
        if(password == ""){
            dispatch(SnackBarActions.showSnackBar("Please input your password"))
            return
        }
        if(confirmPassword == ""){
            dispatch(SnackBarActions.showSnackBar("Please confirm your password"))
            return
        }

        /* check password */
        if(confirmPassword != password){
            dispatch(SnackBarActions.showSnackBar("Password mismatched, please check"))
            return
        }

        /* Show Backdrop*/
        dispatch(BackdropActions.openBackdrop())

        /* Sign up */
        axios.post(`${StateAPI.AuthProtocol}://${StateAPI.AuthHostAddr}:${StateAPI.AuthPort}${StateAPI.AuthBaseURL}/${StateAPI.AuthAPI.SignUp}`, {
            "email": email,
            "password": password
        })
        .then((response) => {
            /* Disabled Backdrop */
            dispatch(BackdropActions.closeBackdrop())

            /* Change state */
            changeSignupState(true)

            /* Show snack bar */
            dispatch(SnackBarActions.showSnackBar(`Successfully Signup, you can signin through ${email}`))
        })
        .catch((error) => {
            if(error.response){
                /* Disabled Backdrop */
                dispatch(BackdropActions.closeBackdrop())

                /* Authentication error */
                if(error.response.status === 500){
                    dispatch(SnackBarActions.showSnackBar("Server internal error, please try agagin later."))
                }
            }else if(error.request){
                /* Disabled Backdrop */
                dispatch(BackdropActions.closeBackdrop())
                
                /* No response */
                dispatch(SnackBarActions.showSnackBar("Can't connect to server, check your network!"))
            }
            return;
        })
    }

    /* Bind callback for button */
    FormFormat.buttons[0].callback = handleSignUpClicked

    if(signupState){
        return <Navigate to="/signin" replace={true} />
    } else{
        return (
            <OuterContainer>
                <PageHeader />
                    <Slide 
                        direction={"right"} 
                        in={true}
                        mountOnEnter 
                        unmountOnExit
                    >
                        <SignupPageContainer>
                            <UserForm FormFormat={FormFormat}/>
                        </SignupPageContainer>
                    </Slide>
                <PageFooter Stick={true} />
            </OuterContainer>
        )
    }
}

export default SignupPage;

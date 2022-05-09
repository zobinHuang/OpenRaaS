import React from 'react';
import { Link } from "react-router-dom";
import axios from 'axios'
import BackgroundImage from '../../Statics/Images/bar_image.svg'
import styled from 'styled-components';
import Avatar from '@mui/material/Avatar';
import UestcLogo from '../../Statics/Images/uestc.png'
import Button from '@mui/material/Button';
import { useDispatch, useSelector } from 'react-redux';
import { actions as AuthActions } from '../../Data/Reducers/authReducer';
import { actions as TerminalsActions } from '../../Data/Reducers/terminalReducer';
import { actions as UserSidebarActions } from '../../Data/Reducers/userSidebarReducer';
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'

const PageHeaderContainer = styled.div`
    width: 100%;
    height: calc(10vh);
    margin: 0px;
    padding: 0px;
    position: fixed;
    top: 0;
    z-index: 1;
    background-color: ${ ({Color}) => Color ? Color : "#0057a3" };
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background-image: url(${BackgroundImage});
    /* background-size: contain; */
`

const TitleLogoContainer = styled.div`
    position: absolute;
    left: 30px;
    margin: 0px;
    padding: 0px;
    display: flex;
    align-items: center;
    justify-content: center;
`

const LogoContainer = styled.img`
    width: 50px;
`

const ButtonGroupContainer = styled.div`
    position: absolute;
    right: 50px;
    width: calc(15vw);
    display: flex;
    align-items: center;
    justify-content: space-evenly;
`

const UserInfoAvatarContainer = styled.div`
    position: absolute;
    right: 50px;
    width: calc(25vw);
    display: flex;
    align-items: center;
    justify-content: space-evenly;
`

const AvatarContainer = styled.div`
`

const UserInfoContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: space-evenly;
`

const UserName = styled.p`
    color: #FFFFFF;
    margin: 0px;
    padding: 0px;
`

const TitleContainer = styled.div`
    margin: 0px 0px 0px 20px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const PageTitle = styled.h1`
    color: #FFFFFF;
    margin: 0px;
    padding: 0px;
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
`

const PageSubTitle = styled.h4`
    color: #FFFFFF;
    margin: 0px;
    padding: 0px;
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
`

const PageHeader = (props) => {
    // get global state
    const StateAuth = useSelector(state => state.auth.StateAuth)
    const StateInfo = useSelector(state => state.info.StateInfo)
    const StateAPI = useSelector(state => state.api.StateAPI)

    const dispatch = useDispatch()

    // handle signout
    const handleSignout = async (event) => {
        event.preventDefault();
        axios.post(`${StateAPI.AuthProtocol}://${StateAPI.AuthHostAddr}:${StateAPI.AuthPort}${StateAPI.AuthBaseURL}/${StateAPI.AuthAPI.SignOut}`, 
            {
                /* Signout request json body, left empty */
            },
            {
                headers: {
                    'Authorization': `Bearer ${StateAuth.idToken}`,
                }
            }
        ).then((response) => {
            /* Change state */
            let un = StateAuth.username
            dispatch(AuthActions.logout())
            dispatch(TerminalsActions.reset())
            dispatch(UserSidebarActions.reset())

            /* Show snack bar */
            dispatch(SnackBarActions.showSnackBar(`Signout from user ${un}`))
        }).catch((error) => {
            if(error.response){
                /* Authentication Out-of-date */
                if(error.response.status === 401){
                    /* Change state */
                    dispatch(AuthActions.logout())
                }
            /* No response */
            }else if(error.request){
                /* Change state */
                dispatch(AuthActions.logout())

                /* Show snack bar */
                dispatch(SnackBarActions.showSnackBar("Can't connect to server, check your network!"))
            }  
        })
    }

    return (
        <PageHeaderContainer Color={StateInfo.ThemeColor}>
            {/* Title, Slogan and Logo */}
            <TitleLogoContainer>
                <LogoContainer src={UestcLogo}/>
                <Link to="/" style={{ textDecoration: 'none' }}>
                    <TitleContainer>
                        <PageTitle>{StateInfo.EnName}</PageTitle>
                        <PageSubTitle>{StateInfo.Slogan}</PageSubTitle>
                    </TitleContainer>
                </Link>
            </TitleLogoContainer>
            
            {/* User Information Area */}
            {StateAuth.isLogin ? (
                    <UserInfoAvatarContainer>
                        <AvatarContainer>
                            <Link to="/user" style={{ textDecoration: 'none' }}>
                                <Avatar src={StateAuth.avatarSrc}/>
                            </Link>
                        </AvatarContainer>
                        <UserInfoContainer>
                            <Link to="/user" style={{ textDecoration: 'none' }}>
                                <UserName>{StateAuth.username}</UserName>
                            </Link>
                        </UserInfoContainer>
                        <Button
                            variant="contained"
                            onClick={handleSignout}
                        >
                            Sign Out
                        </Button>
                    </UserInfoAvatarContainer>
                ) : (
                    <ButtonGroupContainer>
                        <Link to="/signin" style={{ textDecoration: 'none' }}>
                            <Button  variant="contained">Sign In</Button>
                        </Link>
                        <Link to="/signup" style={{ textDecoration: 'none' }}>
                            <Button variant="contained">Sign Up</Button>
                        </Link>
                    </ButtonGroupContainer>
                )
            }
        </PageHeaderContainer>
    )
}

export default PageHeader

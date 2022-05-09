import { React, useEffect, useState } from 'react';
import styled from 'styled-components';
import { Navigate } from "react-router-dom";
import PageHeader from '../../Components/Header/header';
import PageFooter from '../../Components/Footer/footer';
import SideBar from '../../Components/SideBar/sidebar';
import SidebarFormat from "./sidebar_format.json"
import TerminalsPage from './terminals';
import { useDispatch, useSelector } from 'react-redux';
import { actions as UserSidebarActions } from '../../Data/Reducers/userSidebarReducer';

/* Container: Global Container */
const OuterContainer = styled.div`
    width: 100%;
    margin: 0px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding-top: calc(10vh);
`

const UserPageContainer = styled.div`
    /* height: calc(80vh); */
    margin: 0px;
    display: flex;
    flex-direction: row;
    padding-left: calc(8vw);
    min-height: calc(80vh);
`

const SidebarContainer = styled.div`
    margin: 0px;
    padding: 50px 0px;
    width: calc(8vw);
    position: fixed;
    left: 0px;
`

const MainPageContainer = styled.div`
    width: calc(91vw);
   
    padding: 30px 0px;
    display: flex;
    align-items: center;
    justify-content: center;
`

const UserPage = (props) => {
    const dispatch = useDispatch()

    // get global state
    let StateAuth = useSelector(state => state.auth.StateAuth)
    let StateUserSidebar = useSelector(state => state.userSidebar.StateUserSidebar)
    
    // sidebar configuration
    const SidebarConfig = {
        slideDirection: "right",
        sidebarFormat: SidebarFormat,
        showSidebar: StateUserSidebar.showSidebar,
        changeShowSidebar: (value) => dispatch(UserSidebarActions.changeShow(value)),
        currentSelected: StateUserSidebar.currentSelected,
        changeCurrentSelected: (value) => dispatch(UserSidebarActions.changeCurrentSelected(value))
    }

    if(!StateAuth.isLogin){
        return <Navigate to="/" replace={true} />
    } else {
        return (
            <OuterContainer>
                <PageHeader />
                <UserPageContainer>
                    {/* Sidebar Area */}
                    <SidebarContainer>
                        <SideBar SidebarConfig={SidebarConfig} />
                    </SidebarContainer>
                    
                    <MainPageContainer>
                        {
                            /* My Cloud Page */
                            StateUserSidebar.currentSelected == "my_cloud" && 
                            <div>My Cloud Page</div>
                        }

                        {
                            /* Terminals Page */
                            StateUserSidebar.currentSelected == "terminals" && 
                            <TerminalsPage />
                        }

                        {
                            /* Perference Page */
                            StateUserSidebar.currentSelected == "perference" && 
                            <div>Perference Page</div>
                        }
                    </MainPageContainer>
                </UserPageContainer>
                <PageFooter />
            </OuterContainer>
        )
    }
}

export default UserPage;

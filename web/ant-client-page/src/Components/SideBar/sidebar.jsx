import React, { useState } from 'react';
import styled from 'styled-components';
import Slide from '@mui/material/Slide';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';

/* Container: Header */
const SideBarContainer = styled.div`
    width: calc(8vw);
    height: calc(70vh);
    margin: 0px;
    padding: 0px;
    background-color: ${ ({Color}) => Color ? Color : "#e9e9e9" };
    background-size: cover;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const ButtonGroupContainer = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
`

const SideBar = (props) => {

    const {SidebarConfig} = props;

    const HandleClick = async (event) => {
        event.preventDefault();
        SidebarConfig.changeCurrentSelected(event.target.id)
    }

    return (
        <Slide 
            direction={SidebarConfig.slideDirection} 
            in={SidebarConfig.showSidebar}
            mountOnEnter 
            unmountOnExit
        >
            <SideBarContainer>
                <ButtonGroupContainer>
                    <Stack spacing={2} direction="column">
                        {SidebarConfig.sidebarFormat.buttons.map(
                            (button, index) => (
                                SidebarConfig.currentSelected == button.id ?
                                <Button
                                    id={button.id}
                                    key={`sidebar_btn_${index}`} 
                                    size="small"
                                    color="success"
                                    variant="contained"
                                    onClick={HandleClick}
                                    style={{textTransform: 'none'}}
                                >
                                    {button.name}
                                </Button> :
                                <Button
                                    id={button.id}
                                    key={`sidebar_btn_${index}`} 
                                    size="small"
                                    variant="contained"
                                    onClick={HandleClick}
                                    style={{textTransform: 'none'}}
                                >
                                    {button.name}
                                </Button>
                            )
                        )}
                    </Stack>
                </ButtonGroupContainer>
            </SideBarContainer>
        </Slide>
    )
}

export default SideBar
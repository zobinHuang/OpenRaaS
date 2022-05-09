import React from 'react';
import styled from 'styled-components';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';

const UserFormContainer = styled.div`
    width: 500px;
    margin: 0px;
    padding: 30px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px solid ${ ({Color}) => Color ? Color : "#0057a3" };
    border-radius: 20px;
    box-shadow: 5px 5px 5px ${ ({Color}) => Color ? Color : "#0057a3" };
    /* background-color: ${ ({Color}) => Color ? Color : "#0057a3" }; */
`

const UserFormTitleContainer = styled.div`
    width: 100%;
    margin: 30px 0px 0px 0px;
    padding: 0px;
    display: flex;
    align-items: center;
    justify-content: space-around;
`

const UserFormTitle = styled.h2`
    color: ${ ({Color}) => Color ? Color : "#0057a3" };
    margin: 0px;
    padding: 0px;
`

const UserFormEntriesContainer = styled.div`
    width: 100%;
    margin: 0px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const ButtonGroupContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-evenly;
`


const UserForm = (props) => {

    const {FormFormat} = props;
    
    return (
        <UserFormContainer>
            <UserFormTitleContainer>
                <UserFormTitle>{FormFormat.title}</UserFormTitle>
            </UserFormTitleContainer>
            <UserFormEntriesContainer>
                {FormFormat.entries.map(
                    (entry, index) => (
                        !entry.isPassword ?
                        <TextField 
                            fullWidth 
                            margin="normal" 
                            key={`btn_${FormFormat.title}_${entry.name}`} 
                            id={`btn_${FormFormat.title}_${entry.name}`} 
                            label={entry.name} 
                            variant="filled" 
                            onChange={(event) => {
                                entry.changeState(event.target.value)
                            }}
                        /> :
                        <TextField 
                            fullWidth 
                            type="password" 
                            margin="normal" 
                            key={`btn_${FormFormat.title}_${entry.name}`} 
                            id={`btn_${FormFormat.title}_${entry.name}`} 
                            label={entry.name} 
                            variant="filled" 
                            onChange={(event) => {
                                entry.changeState(event.target.value);
                            }}
                        />
                    )
                )}
            </UserFormEntriesContainer>
            <ButtonGroupContainer>
                {FormFormat.buttons.map(
                    (button, index) => (
                        <Button  
                            variant="contained"
                            key={`user_form_btn_${index}`}
                            onClick={button.callback}
                        >
                            {button.name}
                        </Button>
                    )
                )}
            </ButtonGroupContainer>
        </UserFormContainer>
    )
}

export default UserForm
import styled from 'styled-components';

const TagerContainer = styled.div`
    width: max-content;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    margin-right: 20px;
`

const TagerKey = styled.div`
    width: max-content;
    height: inherit;
    background-color: #0057a3;
    color: #ffffff;
    align-items: center;
    justify-content: center;
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
    font-size: small;
    padding: 5px 5px;
    border-radius: 5px 0px 0px 5px;
    box-shadow: 1px 1px 1px #a0a0a0;
` 

const TagerValue = styled.div`
    width: max-content;
    height: inherit;
    background-color: #d9d9d9;
    color: #000000;
    align-items: center;
    justify-content: center;
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
    font-size: small;
    padding: 5px 5px;
    border-radius: 0px 5px 5px 0px;
    box-shadow: 1px 1px 1px #a0a0a0;
` 

const Tager = (props) => {
    const { TagerConfig } = props
    
    return <TagerContainer>
        <TagerKey>
            {TagerConfig.key}
        </TagerKey>
            
        <TagerValue>
            {TagerConfig.value}
        </TagerValue>
    </TagerContainer>
}

export default Tager
import React from 'react';
import styled from 'styled-components';
import BackgroundImage from '../../Statics/Images/bar_image.svg'
import { useDispatch, useSelector } from 'react-redux';

/* Container: Header */
const PageFooterContainer = styled.div`
    width: 100%;
    height: calc(10vh);
    margin: 0px;
    padding: 0px;
    position: ${ ({Stick}) => Stick===true ? "fixed" : "relative" };
    bottom: ${ ({Stick}) => Stick===true ? 0 : null };
    z-index: 1;
    background-image: url(${BackgroundImage});
    background-size: cover;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

`

const FooterInfo = styled.p`
    color: #FFFFFF;
    margin: 0px;
    padding: 0px;
    font-family: 'Gill Sans', 'Gill Sans MT', Calibri, 'Trebuchet MS', sans-serif;
`

const PageFooter = (props) => {

    const { Stick } = props;
    
    const StateInfo = useSelector(state => state.info.StateInfo)

    return (
        <PageFooterContainer Stick={Stick} Color={StateInfo.ThemeColor}>
            <FooterInfo>© {StateInfo.Version} Maintained By {StateInfo.Maintainer}</FooterInfo>
        </PageFooterContainer>
    )
}

export default PageFooter
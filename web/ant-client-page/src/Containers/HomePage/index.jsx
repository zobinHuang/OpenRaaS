import React from 'react';
import styled from 'styled-components';
import VideoPlayer from "react-background-video-player";
import BackgroundVideo from '../../Statics/Video/background_video.mp4'
import PageHeader from '../../Components/Header/header';
import PageFooter from '../../Components/Footer/footer';

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

const HomePageContainer = styled.div`
    width: 100%;
    height: calc(80vh);
    margin: 0px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow-x: hidden;
    overflow-y: hidden;
`

const HomePageSlogan = styled.h1`
    color: #FFFFFF;
    font-size: 50px;
    margin-bottom: 0px;
`

const HomePageSubSlogan = styled.h4`
    color: #FFFFFF;
    margin-top: 0px;
    font-size: 30px;
`

const HomePage = (props) => {
    return (
        <OuterContainer>
            <PageHeader />
            <HomePageContainer>
                <video 
                    autoPlay loop muted
                    style={{
                        width: "100%",
                        position: "absolute",
                        zIndex: "-1"
                    }}
                >
                    <source 
                        src={BackgroundVideo}
                        type="video/mp4"
                    />
                </video>
                <HomePageSlogan>随愿共享，随你所想</HomePageSlogan>
                <HomePageSubSlogan>Share At Your Will</HomePageSubSlogan>
            </HomePageContainer>
            <PageFooter Stick={true} />
        </OuterContainer>
    )
}

export default HomePage;

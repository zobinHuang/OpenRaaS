import styled from 'styled-components';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Button from '@mui/material/Button';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';

const BreadCrumbsContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: ${ ({Position}) => Position ? Position : "left" };
`

const BreadCrumbsInnerContainer = styled.div`
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #e3e3e3;
    border-radius: 0px;
    padding: 10px 30px;
`

const BreadCrumbsNav = (props) => {
    // Format of BreadCrumbsNavConfig
    // {
    //     currentSelectedNavID: "app_market",
    //     handleClick: handleClickNav
    //     navIcons: {
    //         "app_type": <SwipeLeftIcon />,
    //         "app_market": <AppIcon />,
    //     },
    //     navs: [
    //         {
    //             name: "App Type",
    //             id: "app_type",
    //         },
    //         {
    //             name: "App Market",
    //             id: "app_market",
    //         }
    //     ],
    //     position: "left",
    //     disabled: true,
    // }

    const { BreadCrumbsNavConfig } = props;

    return (
    <BreadCrumbsContainer
        Position={BreadCrumbsNavConfig.position && BreadCrumbsNavConfig.position}
    ><BreadCrumbsInnerContainer>
        <Breadcrumbs
            separator={<NavigateNextIcon fontSize='large'/>}
        >
            {
                BreadCrumbsNavConfig.navs.map(
                    (nav, index) => {
                        if(nav.id === BreadCrumbsNavConfig.currentSelectedNavID){
                            return(
                                <Button
                                    key={`breadCurmbsNavButton_${index}`}
                                    variant="outlined"
                                    startIcon={BreadCrumbsNavConfig.navIcons[`${nav.id}`]}
                                    disabled={BreadCrumbsNavConfig.disabled}
                                    onClick={() => {
                                        BreadCrumbsNavConfig.handleClick(nav.id)
                                    }}
                                    sx={ !BreadCrumbsNavConfig.disabled ? {
                                        ":hover": {
                                            backgroundColor:"#840000",
                                            border: "1px solid #840000"
                                        },
                                        fontWeight: 800, 
                                        backgroundColor: "#d10000", 
                                        padding: "0px 10px",
                                        color: "#ffffff",
                                        border: "1px solid #d10000"
                                    } : {
                                        fontWeight: 800, 
                                        backgroundColor: "#41414196", 
                                        padding: "0px 10px",
                                        color: "#ffffff",
                                        border: "1px solid #41414196"
                                    }}
                                >
                                    {nav.name}
                                </Button>
                            )
                        } else {
                            return (
                                <Button 
                                    key={`breadCurmbsNavButton_${index}`}
                                    variant="outlined"
                                    startIcon={BreadCrumbsNavConfig.navIcons[`${nav.id}`]}
                                    onClick={() => {
                                        BreadCrumbsNavConfig.handleClick(nav.id)
                                    }}
                                    disabled={BreadCrumbsNavConfig.disabled}
                                    sx={ !BreadCrumbsNavConfig.disabled ? {
                                        ":hover": {
                                            backgroundColor:"#002749",
                                            border: "1px solid #002749"
                                        },
                                        fontWeight: 800, 
                                        backgroundColor: "#0057a3",
                                        padding: "0px 10px",
                                        color: "#ffffff",
                                        border: "1px solid #0057a3"
                                    } : {
                                        fontWeight: 800, 
                                        backgroundColor: "#41414196",
                                        padding: "0px 10px",
                                        color: "#ffffff",
                                        border: "1px solid #41414196"
                                    }}
                                >
                                    {nav.name}
                                </Button>
                        )}
                    }
                )
            }
        </Breadcrumbs>
    </BreadCrumbsInnerContainer></BreadCrumbsContainer>)
}

export default BreadCrumbsNav
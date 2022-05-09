import React from 'react';
import styled from 'styled-components';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import { useDispatch, useSelector } from 'react-redux';

const StepperContainer = styled.div`
    width: 100%;
    margin: 0px;
    padding: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const CurrentStepContainer = styled.div`
    margin-top: 15px;
    padding: 10px 0px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    background-color: #f3f3f3;
    border-radius: 10px;
`

const RollStepper = (props) => {
    // Format of RollStepperConfig
    // {
    //     steps: [
    //         {
    //             name: "Creating Instance",
    //             id: "creating",
    //             state: "inStep",    // beforeStep, inStep, afterStep, failedStep
    //             descriptionMessage: "Try to create instance",
    //             failedMessage: "Incorrect Configuration Parameters"
    //         },
    //         {
    //             name: "Scheduling Instance",
    //             id: "scheduling",
    //             state: "beforeStep",
    //             descriptionMessage: "Try to schedule instance",
    //             failedMessage: "Failed to Scheduling on Platform"
    //         }
    //     ],
    //  
    //     currentStep: "creating",  
    // }  

    const { RollStepperConfig } = props;

    return (
        <Box sx={{ width: '100%' }}>
            {/* Stepper */}
            <Stepper activeStep={RollStepperConfig.currentStepIndex}>
                {
                    RollStepperConfig.steps.map((step, index) => {
                        const labelProps = {};
                        if(step.state === "failedStep"){
                            labelProps.optional = (
                                <Typography variant="caption" color="error">
                                  {step.failedMessage}
                                </Typography>
                              );
                              labelProps.error = true;
                        }

                        return (
                            <Step key={step.name}>
                              <StepLabel {...labelProps}>{step.name}</StepLabel>
                            </Step>
                        );
                    })
                }
            </Stepper>

            {/* Description Message */}
            <CurrentStepContainer>
                {RollStepperConfig.currentStepIndex === RollStepperConfig.steps.length-1 ? <CheckCircleIcon style={{"margin-right": "20px"}} /> : <CircularProgress style={{"marginRight": "20px"}} size="1.5rem" color="success" />}
                {RollStepperConfig.steps[RollStepperConfig.currentStepIndex].descriptionMessage}
            </CurrentStepContainer>
        </Box>
    )
}

export default RollStepper
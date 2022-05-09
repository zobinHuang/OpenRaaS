import React from 'react';
import styled from 'styled-components';

const LogViewerContainer = styled.div`
    display: flex;
    flex-direction: column;
    background-color: ${ ({Color}) => Color ? Color : "#000000" };
    height: ${ ({Height}) => Height ? Height : "300px" };
    margin: 10px 20px;
    padding: 10px 10px;
    overflow-x: hidden;
    overflow-y: scroll;
`

const LogViewerEntry = styled.div`
    margin: 2px 0px;
    display: flex;
`

const LogViewerLine = styled.p`
    margin: 0px;
    margin-right: 10px;
    color: ${ ({Color}) => Color ? Color : "#ffffff" };
    font-family: "Lucida Console", "Courier New", monospace;
    font-size: 10px;
`

const LogViewerPriority = styled.p`
    margin: 0px;
    margin-right: 10px;
    color: ${ ({Priority}) => {
        if(Priority === "ERROR") return "#ff0000"
        else if(Priority === "SUCCESS") return "#00ffa2"
        else if(Priority === "WARN") return "#fffb00"
        else if(Priority === "INFO") return "#0084ff"
        else return "#ffffff"
    } };
    font-family: "Lucida Console", "Courier New", monospace;
    font-size: 10px;
    font-weight: bolder;
`

const LogViewerTime = styled.p`
    margin: 0px;
    margin-right: 10px;
    color: ${ ({Color}) => Color ? Color : "#ffffff" };
    font-family: "Lucida Console", "Courier New", monospace;
    font-size: 10px;
    font-weight: bolder;
`

const LogViewerContent = styled.p`
    margin: 0px;
    color: ${ ({Color}) => Color ? Color : "#ffffff" };
    font-family: "Lucida Console", "Courier New", monospace;
    font-size: 10px;
    font-weight: bolder;
`

const LogViewer = (props) => {

    // Format of LogViewerConfiguration
    // {
    //     entries: [
    //         {
    //             "priority": "error",
    //             "content": "error message",
    //         }
    //     ],
    // }

    const { LogViewerConfiguration } = props;

    return (
        <LogViewerContainer Height={ LogViewerConfiguration.height && LogViewerConfiguration.height }>
            {LogViewerConfiguration.entries.map(
                (entry, index) => {
                    return <LogViewerEntry key={`log_viewer_entry_${index}`}>
                        <LogViewerLine key={`log_viewer_line_${index}`}>{index+1}</LogViewerLine>
                        <LogViewerTime key={`log_viewer_time_${index}`}>[{entry.time}]</LogViewerTime>
                        <LogViewerPriority 
                            key={`log_viewer_priority_${index}`}
                            Priority={entry.priority}
                        >
                            {entry.priority}
                        </LogViewerPriority>
                        <LogViewerContent key={`log_viewer_content_${index}`}>{entry.content}</LogViewerContent>
                    </LogViewerEntry>
                }
            )}
        </LogViewerContainer>
    )
}

export default LogViewer
import React from 'react';
import styled from 'styled-components';
import Table from '@mui/material/Table';
import Button from '@mui/material/Button';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import TextField from '@mui/material/TextField';
import Skeleton from '@mui/material/Skeleton';
import Paper from '@mui/material/Paper';
import Pagination from '@mui/material/Pagination';
import ArrowDropDownCircleIcon from '@mui/icons-material/ArrowDropDownCircle';
import SearchIcon from '@mui/icons-material/Search';
import WorkOutlineIcon from '@mui/icons-material/WorkOutline';

const RealTimeTableContainer = styled.div`
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
`

const SearchbarContainer = styled.div`
    width: 90%;
    margin: 20px 0px;
    padding: 10px 10px;
    display: flex;
    align-items: center;
    border: 1px solid #a6a6a6;
    border-radius: 10px;
    justify-content: space-around;
`

const SearchbarInnerContainer = styled.div`
    width: 85%;
`

const NoRowsContainer = styled.div`
    width: 100%;
    min-height: 300px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

`

const NoRowsPrompt = styled.h2`
    color: #c4c4c4;
`

const TableHeadTitle = styled.h5`
    color: #000000;
    font-size: 15px;
    margin: 0px;
`

const StyledHeaderCell = styled.div`
    display: flex;
    align-items: center;
    justify-content: center;
`

const StyledTableRow = styled(TableRow)`
    &:hover {
        cursor: ${ ({hoverdisabled}) => hoverdisabled ? "" : "pointer" };;
    }
`

const SelectedTableRow = styled(StyledTableRow)`
    background-color: #ffcaca;
`

const PaginationContainer = styled.div`
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: right;
    margin: 10px 10px;
    padding-right: 20px;
`

export const TABLE_STATE_LOADING = "table_state_loading"
export const TABLE_STATE_LOADED = "table_state_loaded"

const RealTimeTable = (props) => {

    // format of RealTimeTableConfig
    // {
    //     "heads": [
    //         {
    //             "index": "application",
    //             "name": "Application",
    //             "align": "right",
    //         },
    //         {
    //             "index": "update_time",
    //             "name": "Update Time",
    //             "align": "right",
    //         },
    //         ...
    //     ],
    //
    //     "rows": [
    //         {
    //             "value": {
    //                 "application": "road rash",
    //                 "update_time": "Apr.21 2022"
    //             }
    //         }
    //     ],
    //
    //     "tableState": "loaded" // loaded & loading
    // }

    const { RealTimeTableConfig } = props;

    return <RealTimeTableContainer>
            <SearchbarContainer>
                <SearchbarInnerContainer>
                    <TextField 
                        fullWidth
                        size="small"
                        sx={{zIndex: "0"}}
                        label={RealTimeTableConfig.searchBarPlaceHolder} 
                        id="fullWidth"
                        variant='outlined'
                        color="secondary"
                        disabled={RealTimeTableConfig.disabled}
                    />
                </SearchbarInnerContainer>

                <Button 
                    variant="contained" 
                    endIcon={<SearchIcon />}
                    disabled={RealTimeTableConfig.disabled}
                >
                Search
                </Button>
            </SearchbarContainer>

            <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">

                {/* ----------------------- Table Heads ----------------------- */}
                <TableHead>
                    <TableRow>
                        {
                            RealTimeTableConfig.heads.map(
                                (head, index) => {
                                    return (
                                        <TableCell 
                                            align="center"
                                            key={`readtimeTableHeadCell_${index}`}
                                        >
                                            <StyledHeaderCell>
                                                {head.canOrder && 
                                                    <IconButton
                                                        disabled={RealTimeTableConfig.disabled}
                                                        onClick={head.OrderCallback}
                                                    > 
                                                        <ArrowDropDownCircleIcon /> 
                                                    </IconButton>
                                                }
                                                <TableHeadTitle>{head.name}</TableHeadTitle>
                                            </StyledHeaderCell>
                                        </TableCell>
                                    )
                                }
                            )
                        }
                    </TableRow>
                </TableHead>

                {/* ----------------------- Table Body ----------------------- */}
                { /* Case: table is loaded and length is larger than 0 */}
                { (RealTimeTableConfig.rows.length > 0 && RealTimeTableConfig.tableState === TABLE_STATE_LOADED) &&
                    <TableBody>
                    {
                        RealTimeTableConfig.rows.map(
                            (row, row_index) => {
                                if(row.selected){
                                    return (
                                        <SelectedTableRow
                                            key={`realtimeTableRow_${row_index}`}
                                            onClick={!RealTimeTableConfig.disabled && RealTimeTableConfig.handleClickOnRow}
                                            id={`${row_index}`}
                                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                            hoverdisabled={RealTimeTableConfig.disabled}
                                        >
                                            {
                                                Object.keys(row.values).map(
                                                    (value_key, column_index) => {
                                                        return (
                                                            <TableCell
                                                                key={`realtimeTable_Row${row_index}_Cell${column_index}`}
                                                                id={`${row_index}`} 
                                                                align="center"
                                                            >
                                                                {row.values[value_key]}
                                                            </TableCell>
                                                        )
                                                    }
                                                )
                                            }
                                        </SelectedTableRow>
                                    )
                                } else {
                                    return (
                                        <StyledTableRow
                                            key={`realtimeTableRow_${row_index}`}
                                            hover={!RealTimeTableConfig.disabled}
                                            hoverdisabled={RealTimeTableConfig.disabled}
                                            onClick={!RealTimeTableConfig.disabled && RealTimeTableConfig.handleClickOnRow}
                                            id={`${row_index}`}
                                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                        >
                                            {
                                                Object.keys(row.values).map(
                                                    (value_key, column_index) => {
                                                        return (
                                                            <TableCell
                                                                key={`realtimeTable_Row${row_index}_Cell${column_index}`}
                                                                id={`${row_index}`} align="center">
                                                                {row.values[value_key]}
                                                            </TableCell>
                                                        )
                                                    }
                                                )
                                            }
                                        </StyledTableRow>
                                    ) 
                                }
                            }
                        )
                    }
                    </TableBody>
                }

                { /* Case: table is still loading */}
                { RealTimeTableConfig.tableState === TABLE_STATE_LOADING && 
                    <TableBody>
                        <TableRow>
                            {
                                RealTimeTableConfig.heads.map(
                                    (head, skeleton_index) => {
                                        return <TableCell
                                            key={`skeleton_${skeleton_index}`}
                                            id={`skeleton_${skeleton_index}`} 
                                            align="center"
                                        >
                                            <Skeleton animation="wave" />
                                        </TableCell>
                                    }
                                )
                            }
                        </TableRow>
                        <TableRow>
                            {
                                RealTimeTableConfig.heads.map(
                                    (head, skeleton_index) => {
                                        return <TableCell
                                            key={`skeleton_${skeleton_index}`}
                                            id={`skeleton_${skeleton_index}`} 
                                            align="center"
                                        >
                                            <Skeleton animation="wave" />
                                        </TableCell>
                                    }
                                )
                            }
                        </TableRow>
                        <TableRow>
                            {
                                RealTimeTableConfig.heads.map(
                                    (head, skeleton_index) => {
                                        return <TableCell
                                            key={`skeleton_${skeleton_index}`}
                                            id={`skeleton_${skeleton_index}`} 
                                            align="center"
                                        >
                                            <Skeleton animation="wave" />
                                        </TableCell>
                                    }
                                )
                            }
                        </TableRow>
                    </TableBody> 
                }
            </Table>
        </TableContainer>
        
        { /* Case: table is loaded yet length is equals to 0 */}
        { (RealTimeTableConfig.rows.length === 0 && RealTimeTableConfig.tableState === TABLE_STATE_LOADED) && <NoRowsContainer>
            <WorkOutlineIcon color="disabled" sx={{ fontSize: 150, margin: 0 }}/>
            <NoRowsPrompt>No Results</NoRowsPrompt>
        </NoRowsContainer>}
    
        <PaginationContainer>
            {RealTimeTableConfig.rows.length > 0 && <Pagination 
                count={RealTimeTableConfig.rowAmount % RealTimeTableConfig.maxRowsPerPage === 0 ? 
                    Math.trunc(RealTimeTableConfig.rowAmount / RealTimeTableConfig.maxRowsPerPage) :
                    (Math.trunc(RealTimeTableConfig.rowAmount / RealTimeTableConfig.maxRowsPerPage)+1)
                }
                onChange={RealTimeTableConfig.handleChangePagination}
                page={RealTimeTableConfig.currentSelectedPageIndex}
                variant="outlined"
                disabled={RealTimeTableConfig.disabled}
                shape="rounded" 
            />}
        </PaginationContainer>
    </RealTimeTableContainer>
}

export default RealTimeTable
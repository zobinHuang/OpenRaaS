import * as React from 'react';
import Snackbar from '@mui/material/Snackbar';
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';
import { useSelector, useDispatch } from 'react-redux'
import { actions as SnackBarActions } from '../../Data/Reducers/snackBarReducer'

const SnackBar = (props) => {
    const dispatch = useDispatch()

    const StateSnackBar = useSelector(state => state.snackbar.StateSnackBar)

    const handleClose = () => {
        dispatch(SnackBarActions.closeSnackBar())
    }

    const action = (
        <React.Fragment>
            <IconButton
                size="small"
                aria-label="close"
                color="inherit"
                onClick={handleClose}
            >
            <CloseIcon fontSize="small" />
            </IconButton>
        </React.Fragment>
    );

    return (<Snackbar
        open={StateSnackBar.snackBarEnabled}
        autoHideDuration={6000}
        onClose={handleClose}
        message={StateSnackBar.snackBarContent}
        action={action}
    />)
}

export default SnackBar;
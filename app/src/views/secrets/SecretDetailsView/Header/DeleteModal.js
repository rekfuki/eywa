import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
  Box,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Input,
  Typography,
  makeStyles
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(3)
  },
  fontWeightBold: {
    fontWeight: theme.typography.fontWeightBold
  },
  colorRed: {
    color: theme.palette.error.main,
    fontWeight: theme.typography.fontWeightBold
  },
  largeTitle: {
    "font-size": 20
  }
}));

const DeleteModal = ({
  secretName,
  onDelete,
  onClose,
  open
}) => {
  const classes = useStyles();
  const [deleteEnabled, setDeleteEnabled] = useState(false);

  const handleChange = (event) => {
    event.persist();
    setDeleteEnabled(event.target.value === secretName)
  };

  return (
    <Dialog open={open} onClose={onClose} aria-labelledby="form-dialog-title">
      <DialogTitle id="form-dialog-title"></DialogTitle>
      <DialogContent>
        <Typography variant="h3" component="h3" gutterBottom>
          {`Delete secret ${secretName}?`}
        </Typography>
        <Box mt={3}></Box>
        <Typography>
          <span>
            This action cannot be
          </span>
          {' '}
          <span className={classes.fontWeightBold}>
            undone.
          </span>
          {' '}
          <span>
            This will
          </span>
          {' '}
          <span className={classes.colorRed}>
            permamently delete
          </span>
          {' '}
          <span>
            the secret.
          </span>
        </Typography>
        <Box mt={3}></Box>
        <Typography>
          <span>
            Please enter
          </span>
          {' '}
          <span className={classes.fontWeightBold}>
            {secretName}
          </span>
          {' '}
          <span>
            to confirm.
          </span>
        </Typography>
        <Input fullWidth autoFocus required onChange={handleChange} placeholder={secretName} />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          Cancel
        </Button>
        <Button disabled={!deleteEnabled} onClick={onDelete} color="primary">
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DeleteModal;

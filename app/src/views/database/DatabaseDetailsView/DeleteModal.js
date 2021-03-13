import React, { useState } from 'react';
import {
  Box,
  Button,
  Dialog,
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
  collection,
  onDelete,
  onClose,
  open
}) => {
  const classes = useStyles();
  const [deleteEnabled, setDeleteEnabled] = useState(false);

  if (!collection) {
    return null;
  }

  const collectionName = collection.namespace.substring(collection.namespace.indexOf('.') + 1)

  const handleChange = (event) => {
    event.persist();
    setDeleteEnabled(event.target.value === collectionName)
  };

  const handleDelete = () => {
    onDelete();
  };

  return (
    <Dialog open={open} onClose={onClose} aria-labelledby="form-dialog-title">
      <DialogContent>
        <Typography variant="h3" component="h3" gutterBottom>
          {`Delete collection ${collectionName}?`}
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
            the collection.
          </span>
        </Typography>
        <Box mt={3}></Box>
        <Typography>
          <span>
            Please enter
          </span>
          {' '}
          <span className={classes.fontWeightBold}>
            {collectionName}
          </span>
          {' '}
          <span>
            to confirm.
          </span>
        </Typography>
        <Input fullWidth autoFocus required onChange={handleChange} placeholder={collectionName} />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          Cancel
        </Button>
        <Button disabled={!deleteEnabled} onClick={handleDelete} color="primary">
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DeleteModal;

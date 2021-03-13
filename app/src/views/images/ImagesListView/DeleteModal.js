import React, { useState } from 'react';
import PropTypes from 'prop-types';
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
  imageName,
  className,
  onDelete,
  onClose,
  open,
  ...rest
}) => {
  const classes = useStyles();
  const [deleteEnabled, setDeleteEnabled] = useState(false);

  const handleChange = (event) => {
    event.persist();
    setDeleteEnabled(event.target.value === imageName)
  };

  const handleDelete = () => {
    onDelete();
  };

  return (
    <Dialog open={open} onClose={onClose} aria-labelledby="form-dialog-title">
      <DialogContent>
        <Typography variant="h3" component="h3" gutterBottom>
          {`Delete image ${imageName}?`}
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
            the image.
          </span>
        </Typography>
        <Box mt={3}></Box>
        <Typography>
          <span>
            Please enter
          </span>
          {' '}
          <span className={classes.fontWeightBold}>
            {imageName}
          </span>
          {' '}
          <span>
            to confirm.
          </span>
        </Typography>
        <Input fullWidth autoFocus required onChange={handleChange} placeholder={imageName}/>
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

DeleteModal.propTypes = {
  imageName: PropTypes.string,
  className: PropTypes.string,
  onApply: PropTypes.func,
  onClose: PropTypes.func,
  open: PropTypes.bool.isRequired
};

DeleteModal.defaultProps = {
  onApply: () => {},
  onClose: () => {}
};

export default DeleteModal;

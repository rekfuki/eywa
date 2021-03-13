import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Breadcrumbs,
  Button,
  Grid,
  Link,
  SvgIcon,
  Typography,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import DeleteIcon from '@material-ui/icons/DeleteOutline';
import {
  Edit as EditIcon
} from 'react-feather';
import DeleteModal from './DeleteModal'

const useStyles = makeStyles((theme) => ({
  root: {},
  deleteAction: {
    color: theme.palette.common.white,
    backgroundColor: theme.palette.error.main,
    '&:hover': {
      backgroundColor: theme.palette.error.dark
    }
  }
}));

const Header = ({
  className,
  secret,
  onDeleteSecret,
  ...rest }) => {
  const classes = useStyles();
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);

  const handleDeleteModalOpen = () => {
    setDeleteModalOpen(true);
  };

  const handleDeleteModalClose = () => {
    setDeleteModalOpen(false);
  };

  const handleDeleteModalExecute = () => {
    setDeleteModalOpen(false);
    onDeleteSecret()
  }

  return (
    <Grid
      container
      spacing={3}
      justify="space-between"
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Grid item>
        <Breadcrumbs
          separator={<NavigateNextIcon fontSize="small" />}
          aria-label="breadcrumb"
        >
          <Link
            variant="body1"
            color="inherit"
            to="/app"
            component={RouterLink}
          >
            Dashboard
          </Link>
          <Link
            variant="body1"
            color="inherit"
            to="/app/secrets"
            component={RouterLink}
          >
            Secrets
          </Link>
          <Typography
            variant="body1"
            color="textPrimary"
          >
            {secret.id}
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          {secret.name}
        </Typography>
      </Grid>
      <Grid item>
        <Grid container spacing={3} justify="center">
          <Grid item>
            <Button
              className={classes.deleteAction}
              variant="contained"
              onClick={handleDeleteModalOpen}
              startIcon={
                <SvgIcon fontSize="small">
                  <DeleteIcon />
                </SvgIcon>
              }
            >
              Delete
            </Button>
            <DeleteModal
              secretName={secret.name}
              onDelete={handleDeleteModalExecute}
              onClose={handleDeleteModalClose}
              open={isDeleteModalOpen}
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string
};

export default Header;

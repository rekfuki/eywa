import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Button,
  Breadcrumbs,
  Grid,
  Link,
  Typography,
  SvgIcon,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import {
  PlusCircle as PlusCircleIcon
} from 'react-feather';
import CreateModal from './CreateModal'

const useStyles = makeStyles((theme) => ({
  root: {},
  action: {
    marginBottom: theme.spacing(1),
    '& + &': {
      marginLeft: theme.spacing(1)
    }
  }
}));

const Header = ({
  className,
  refresh,
  ...rest }) => {
  const classes = useStyles();
  const [isCreateModalOpen, setCreateModalOpen] = useState(false);


  const handleCreateModalOpen = () => {
    setCreateModalOpen(true);
  };

  const handleCreateModalClose = (created) => {
    if (created) {
      refresh();
    }

    setCreateModalOpen(false);
  };

  return (
    <Grid
      className={clsx(classes.root, className)}
      container
      justify="space-between"
      spacing={3}
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
          <Typography
            variant="body1"
            color="textPrimary"
          >
            Tokens
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          All Tokens
        </Typography>
      </Grid>
      <Grid item>
        <Grid container spacing={3} justify="center">
          <Grid item>
            <Button
              color="secondary"
              variant="contained"
              onClick={handleCreateModalOpen}
              startIcon={
                <SvgIcon fontSize="small">
                  <PlusCircleIcon />
                </SvgIcon>
              }
            >
              Create Token
            </Button>
            <CreateModal
              onClose={handleCreateModalClose}
              open={isCreateModalOpen}
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

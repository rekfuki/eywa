import React from 'react';
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
  PlusCircle as PlusCircleIcon,
} from 'react-feather';

const useStyles = makeStyles((theme) => ({
  root: {},
  action: {
    marginBottom: theme.spacing(1),
    '& + &': {
      marginLeft: theme.spacing(1)
    }
  }
}));

const Header = ({ className, ...rest }) => {
  const classes = useStyles();

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
            Images
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          All Images
        </Typography>
      </Grid>
      <Grid item>
        <Button
          color="secondary"
          variant="contained"
          component={RouterLink}
          to="/app/images/create"
          startIcon={
            <SvgIcon fontSize="small">
              <PlusCircleIcon />
            </SvgIcon>
          }
        >
          Create Image
        </Button>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string
};

export default Header;

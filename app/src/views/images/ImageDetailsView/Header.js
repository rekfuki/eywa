import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Breadcrumbs,
  Grid,
  Link,
  Typography,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';

const useStyles = makeStyles((theme) => ({
  root: {},
  action: {
    marginBottom: theme.spacing(1),
    '& + &': {
      marginLeft: theme.spacing(1)
    }
  }
}));

const Header = ({ className, imageId, ...rest }) => {
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
          <Link
            variant="body1"
            color="inherit"
            to="/app/images"
            component={RouterLink}
          >
            Images
          </Link>
          <Typography
            variant="body1"
            color="textPrimary"
          >
            {imageId}
          </Typography>
          <Typography
            variant="body1"
            color="textPrimary"
          >
            Build Logs
          </Typography>
        </Breadcrumbs>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string,
  imageId: PropTypes.string.isRequired
};

export default Header;

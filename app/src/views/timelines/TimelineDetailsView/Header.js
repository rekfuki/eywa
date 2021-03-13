import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Box,
  Breadcrumbs,
  Button,
  IconButton,
  Grid,
  Link,
  List,
  ListItem,
  SvgIcon,
  Typography,
  TextField,
  Tooltip,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import useIsMountedRef from 'src/hooks/useIsMountedRef';

const useStyles = makeStyles((theme) => ({
  root: {},
}));

const Header = ({ className, requestId, ...rest }) => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();

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
            to="/app/timelines"
            component={RouterLink}
          >
            Executions
          </Link>
          <Typography
            variant="body1"
            color="textPrimary"
          >
            {requestId}
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          Execution Details
        </Typography>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string,
  requestId: PropTypes.string.isRequired
};

export default Header;

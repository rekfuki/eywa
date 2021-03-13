import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  AppBar,
  Box,
  Hidden,
  IconButton,
  Typography,
  Toolbar,
  makeStyles,
  SvgIcon
} from '@material-ui/core';
import {
  Menu as MenuIcon,
  LogOut as LogOutIcon
} from 'react-feather';
import useAuth from 'src/hooks/useAuth';
import Logo from 'src/components/Logo';
import { THEMES } from 'src/constants';

const useStyles = makeStyles((theme) => ({
  root: {
    zIndex: theme.zIndex.drawer + 100,
    ...theme.name === THEMES.LIGHT ? {
      boxShadow: 'none',
      backgroundColor: theme.palette.primary.main
    } : {},
    ...theme.name === THEMES.ONE_DARK ? {
      backgroundColor: theme.palette.background.default
    } : {}
  },
  toolbar: {
    minHeight: 52
  }
}));

const TopBar = ({
  className,
  onMobileNavOpen,
  ...rest
}) => {
  const classes = useStyles();
  const { user, logout } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Unable to logout', {
        variant: 'error'
      });
    }
  };

  return (
    <AppBar
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Toolbar className={classes.toolbar}>
        <Hidden lgUp>
          <IconButton
            color="inherit"
            onClick={onMobileNavOpen}
          >
            <SvgIcon fontSize="small">
              <MenuIcon />
            </SvgIcon>
          </IconButton>
        </Hidden>
        <Hidden mdDown>
          <RouterLink to="/" style={{ display: "flex", alignItems: "center", textDecoration: "none" }}>
            <Logo width={50} />
            <Typography style={{ color: "white", letterSpacing: "3px", marginLeft: "20px" }}>
              EYWA
            </Typography>
          </RouterLink>
        </Hidden>
        <Box
          ml={2}
          flexGrow={1}
        />
        <IconButton color="inherit" size="small" onClick={handleLogout}>
          <SvgIcon fontSize="small">
            <LogOutIcon />
          </SvgIcon>
          <Typography variant="button" style={{ marginLeft: "10px" }}>
            LOG OUT
          </Typography>
        </IconButton>
      </Toolbar>
    </AppBar>
  );
};

TopBar.propTypes = {
  className: PropTypes.string,
  onMobileNavOpen: PropTypes.func
};

TopBar.defaultProps = {
  onMobileNavOpen: () => { }
};

export default TopBar;

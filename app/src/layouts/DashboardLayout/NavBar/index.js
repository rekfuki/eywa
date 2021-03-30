/* eslint-disable no-use-before-define */
import React, { useEffect } from 'react';
import { useLocation, matchPath } from 'react-router-dom';
import { Link as RouterLink } from 'react-router-dom';
import PerfectScrollbar from 'react-perfect-scrollbar';
import PropTypes from 'prop-types';
import {
  Avatar,
  Box,
  Divider,
  Drawer,
  Hidden,
  List,
  ListSubheader,
  Typography,
  makeStyles
} from '@material-ui/core';
import {
  Lock as LockIcon,
  PieChart as PieChartIcon,
  Command as CommandIcon,
  Disc as DiscIcon,
  Key as KeyIcon,
  Database as DatabaseIcon,
  BookOpen as OpenBookIcon
} from 'react-feather';
import TimelineIcon from '@material-ui/icons/Timeline';
import ViewHeadlineIcon from '@material-ui/icons/ViewHeadline';
import Logo from 'src/components/Logo';
import useAuth from 'src/hooks/useAuth';
import NavItem from './NavItem';

const sections = [
  {
    // subheader: 'Reports',
    items: [
      {
        title: 'Dashboard',
        icon: PieChartIcon,
        href: '/app/dashboard'
      }
    ]
  },
  {
    href: "/app/timelines",
    items: [
      {
        title: "Executions",
        icon: TimelineIcon,
        href: "/app/timelines"
      }
    ]
  },
  {
    href: "/app/logs",
    items: [
      {
        title: "Logs",
        icon: ViewHeadlineIcon,
        href: "/app/logs"
      }
    ]
  },
  {
    href: "/app/images",
    items: [
      {
        title: "Images",
        icon: DiscIcon,
        href: "/app/images"
      }
    ]
  },
  {
    href: "/app/functions",
    items: [
      {
        title: "Functions",
        icon: CommandIcon,
        href: "/app/functions"
      }
    ]
  },
  {
    href: "/app/secrets",
    items: [
      {
        title: "Secrets",
        icon: LockIcon,
        href: "/app/secrets"
      }
    ]
  },
  {
    href: "/app/tokens",
    items: [
      {
        title: "Tokens",
        icon: KeyIcon,
        href: "/app/tokens"
      }
    ]
  },
  {
    href: "/app/timelines",
    items: [
      {
        title: 'Database',
        icon: DatabaseIcon,
        href: '/app/database'
      }
    ]
  }
];

function renderNavItems({
  items,
  pathname,
  depth = 0
}) {
  return (
    <List disablePadding>
      {items.reduce(
        (acc, item) => reduceChildRoutes({ acc, item, pathname, depth }),
        []
      )}
    </List>
  );
}

function reduceChildRoutes({
  acc,
  pathname,
  item,
  depth
}) {
  const key = item.title + depth;

  if (item.items) {
    const open = matchPath(pathname, {
      path: item.href,
      exact: false
    });

    acc.push(
      <NavItem
        raw={item.items}
        depth={depth}
        icon={item.icon}
        info={item.info}
        key={key}
        open={Boolean(open)}
        title={item.title}
      >
        {renderNavItems({
          depth: depth + 1,
          pathname,
          items: item.items
        })}
      </NavItem>
    );
  } else {
    acc.push(
      <NavItem
        raw={item.raw}
        depth={depth}
        href={item.href}
        icon={item.icon}
        info={item.info}
        key={key}
        title={item.title}
      />
    );
  }

  return acc;
}

const useStyles = makeStyles(() => ({
  mobileDrawer: {
    width: 256
  },
  desktopDrawer: {
    width: 256,
    top: 64,
    height: 'calc(100% - 64px)'
  },
  avatar: {
    cursor: 'pointer',
    width: 64,
    height: 64
  }
}));

const NavBar = ({ onMobileClose, openMobile }) => {
  const classes = useStyles();
  const location = useLocation();
  const { user } = useAuth();

  useEffect(() => {
    if (openMobile && onMobileClose) {
      onMobileClose();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname]);

  const content = (
    <Box
      height="100%"
      display="flex"
      flexDirection="column"
    >
      <PerfectScrollbar options={{ suppressScrollX: true }}>
        <Hidden lgUp>
          <Box
            p={2}
            display="flex"
            justifyContent="center"
          >
            <RouterLink to="/">
              <Logo />
            </RouterLink>
          </Box>
        </Hidden>
        <Box p={2}>
          <Box
            display="flex"
            justifyContent="center"
          >
            <Avatar
              alt="User"
              className={classes.avatar}
              src={user.avatar_url}
            />
          </Box>
          <Box
            mt={2}
            textAlign="center"
          >
            <Typography
              variant="h5"
              color="textPrimary"
            >
              {user.oauth_provider_login}
            </Typography>
          </Box>
          <Box
            mt={1}
            textAlign="center"
          >
            <Typography
              variant="h6"
              color="textPrimary"
            >
              {user.name}
            </Typography>
          </Box>
        </Box>
        <Divider />
        <Box p={2}>
          {sections.map((section, index) => (
            <List
              key={index}
              subheader={(
                <ListSubheader
                  disableGutters
                  disableSticky
                >
                  {section.subheader}
                </ListSubheader>
              )}
            >
              {renderNavItems({
                items: section.items,
                pathname: location.pathname
              })}
            </List>
          ))}
        </Box>
        <Divider />
        <Box p={2}>
          {renderNavItems({
            items: [{
              title: 'User Guide',
              href: '/docs'
            }],
            pathname: location.pathname
          })}
          {renderNavItems({
            items: [{
              title: 'REST API docs',
              href: '/api-docs',
              raw: true
            }]
          })}
        </Box>
      </PerfectScrollbar>
    </Box>
  );

  return (
    <>
      <Hidden lgUp>
        <Drawer
          anchor="left"
          classes={{ paper: classes.mobileDrawer }}
          onClose={onMobileClose}
          open={openMobile}
          variant="temporary"
        >
          {content}
        </Drawer>
      </Hidden>
      <Hidden mdDown>
        <Drawer
          anchor="left"
          classes={{ paper: classes.desktopDrawer }}
          open
          variant="persistent"
        >
          {content}
        </Drawer>
      </Hidden>
    </>
  );
};

NavBar.propTypes = {
  onMobileClose: PropTypes.func,
  openMobile: PropTypes.bool
};

export default NavBar;

/* eslint-disable no-use-before-define */
import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import {
  Box,
  Drawer,
  Hidden,
  List,
  makeStyles
} from '@material-ui/core';
import Logo from 'src/components/Logo';
import NavItem from './NavItem.bak';
import NavSection from './NavSection';

const sections = [
  {
    title: 'Overview',
    items: [
      {
        title: 'Welcome',
        path: '/docs/overview/welcome'
      },
      {
        title: 'Getting Started',
        path: '/docs/overview/getting-started'
      }
    ]
  },
  {
    title: 'Executions',
    items: [
      {
        title: 'Overview',
        path: '/docs/executions/overview'
      }
    ]
  },
  {
    title: 'Logs',
    items: [
      {
        title: 'Overview',
        path: '/docs/logs/overview'
      }
    ]
  },
  {
    title: 'Images',
    items: [
      {
        title: 'Create Image',
        path: '/docs/images/create'
      },
      {
        title: 'Go runtime',
        path: '/docs/images/go'
      },
      {
        title: 'NodeJS 14 runtime',
        path: '/docs/images/node14'
      },
      {
        title: 'Python3 runtime',
        path: '/docs/images/python3'
      },
      {
        title: 'Ruby runtime',
        path: '/docs/images/ruby'
      },
      {
        title: 'Custom runtime',
        path: '/docs/images/custom'
      }
    ]
  },
  {
    title: 'Functions',
    items: [
      {
        title: 'Manage',
        path: '/docs/functions/manage'
      },
      {
        title: 'Create',
        path: '/docs/functions/create'
      },
      {
        title: 'Update',
        path: '/docs/functions/update'
      },
      {
        title: 'Delete',
        path: '/docs/functions/delete'
      },
      {
        title: 'Metrics',
        path: '/docs/functions/metrics'
      }
    ]
  },
  {
    title: 'Secrets',
    items: [
      {
        title: 'Create Secrets',
        path: '/docs/secrets/create'
      }
    ]
  },
  {
    title: 'Tokens',
    items: [
      {
        title: 'Create Token',
        path: '/docs/access_tokens/create'
      },
      {
        title: 'Delete Token',
        path: '/docs/access_tokens/delete'
      }
    ]
  }
];

function renderNavItems({ items, depth = 0 }) {
  return (
    <List disablePadding>
      {items.reduce(
        (acc, item) => reduceChildRoutes({ acc, item, depth }),
        []
      )}
    </List>
  );
}

function reduceChildRoutes({
  acc,
  item,
  depth = 0
}) {
  if (item.items) {
    acc.push(
      <NavItem
        depth={depth}
        key={item.href}
        title={item.title}
      >
        {renderNavItems({
          items: item.items,
          depth: depth + 1
        })}
      </NavItem>
    );
  } else {
    acc.push(
      <NavItem
        depth={depth}
        href={item.href}
        key={item.href}
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
  }
}));

const NavBar = ({ onMobileClose, openMobile }) => {
  const classes = useStyles();
  const location = useLocation();

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
      <Hidden lgUp>
        <Box p={2}>
          <RouterLink to="/">
            <Logo />
          </RouterLink>
        </Box>
      </Hidden>
      <Box p={2}>
        {sections.map((section) => (
          <NavSection
            key={section.title}
            pathname={location.pathname}
            {...section}
          />

        ))}
        {/* {renderNavItems({ items })} */}
      </Box>
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

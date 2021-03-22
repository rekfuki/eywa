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
      },
      {
        title: 'Dependencies',
        path: '/docs/overview/dependencies'
      },
      {
        title: 'Environment Variables',
        path: '/docs/overview/environment-variables'
      },
      {
        title: 'Theming',
        path: '/docs/overview/theming'
      },
      {
        title: 'Redux',
        path: '/docs/overview/redux'
      },
      {
        title: 'Server Calls',
        path: '/docs/overview/server-calls'
      },
      {
        title: 'Settings',
        path: '/docs/overview/settings'
      },
      {
        title: 'RTL',
        path: '/docs/overview/rtl'
      },
      {
        title: 'Internationalization',
        path: '/docs/overview/internationalization'
      },
      {
        title: 'Deployment',
        path: '/docs/overview/deployment'
      },
      {
        title: 'Migrating to Next.js',
        path: '/docs/overview/migrating-to-nextjs'
      }
    ]
  },
  {
    title: 'Functions',
    items: [
      {
        title: 'Overview',
        path: '/docs/functions/overview'
      },
      {
        title: 'Code Splitting',
        path: '/docs/routing/code-splitting'
      }
    ]
  },
  {
    title: 'Images',
    items: [
      {
        title: 'Create Image',
        path: '/docs/images/create'
      }
    ]
  },
  {
    title: 'Secrets',
    items: [
      {
        title: 'Create Secrets',
        path: '/docs/secrets/create'
      },
      {
        title: 'Managing Secrets',
        path: '/docs/routing/code-splitting'
      }
    ]
  },
  {
    title: 'Authentication',
    items: [
      {
        title: 'Amplify',
        path: '/docs/authentication/amplify'
      },
      {
        title: 'Auth0',
        path: '/docs/authentication/auth0'
      },
      {
        title: 'Firebase',
        path: '/docs/authentication/firebase'
      },
      {
        title: 'JWT',
        path: '/docs/authentication/jwt'
      }
    ]
  },
  {
    title: 'Guards',
    items: [
      {
        title: 'Guest Guard',
        path: '/docs/guards/guest-guard'
      },
      {
        title: 'Auth Guard',
        path: '/docs/guards/auth-guard'
      },
      {
        title: 'Role Based Guard',
        path: '/docs/guards/role-based-guard'
      }
    ]
  },
  {
    title: 'Analytics',
    items: [
      {
        title: 'Introduction',
        path: '/docs/analytics/introduction'
      },
      {
        title: 'Configuration',
        path: '/docs/analytics/configuration'
      },
      {
        title: 'Event Tracking',
        path: '/docs/analytics/event-tracking'
      }
    ]
  },
  {
    title: 'Support',
    items: [
      {
        title: 'Changelog',
        path: '/docs/support/changelog'
      },
      {
        title: 'Contact',
        path: '/docs/support/contact'
      },
      {
        title: 'Further Support',
        path: '/docs/support/further-support'
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

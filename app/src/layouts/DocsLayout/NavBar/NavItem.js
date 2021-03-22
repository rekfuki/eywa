import { useState } from 'react';
import { NavLink as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Box, Button, Collapse, ListItem, useTheme } from '@material-ui/core';
// import ChevronDownIcon from '../icons/ChevronDown';
// import ChevronRightIcon from '../icons/ChevronRight';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

const NavItem = (props) => {
  const theme = useTheme();
  const { active, children, depth, icon, info, open: openProp, path, title, ...other } = props;
  const [open, setOpen] = useState(openProp);

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  let paddingLeft = 16;

  if (depth > 0) {
    paddingLeft = 32 + 8 * depth;
  }

  // Branch
  if (children) {
    return (
      <ListItem
        disableGutters
        style={{
          display: 'block',
          py: 0
        }}
        {...other}
      >
        <Button
          endIcon={!open ? <KeyboardArrowUpIcon fontSize="small" />
            : <KeyboardArrowDownIcon fontSize="small" />}
          onClick={handleToggle}
          startIcon={icon}
          style={{
            color: 'text.secondary',
            fontWeight: 'medium',
            justifyContent: 'flex-start',
            paddingLeft: `${paddingLeft}px`,
            paddingRight: '8px',
            py: '12px',
            textAlign: 'left',
            textTransform: 'none',
            width: '100%'
          }}
          variant="text"
        >
          <Box styles={{ flexGrow: 1 }}>
            {title}
          </Box>
          {info}
        </Button>
        <Collapse in={open}>
          {children}
        </Collapse>
      </ListItem>
    );
  }

  // Leaf
  return (
    <ListItem
      disableGutters
      style={{
        display: 'flex',
        py: 0
      }}
    >
      <Button
        component={path && RouterLink}
        startIcon={icon}
        style={{
          color: theme.palette.text.primary,
          fontWeight: 'fontWeightMedium',
          justifyContent: 'flex-start',
          textAlign: 'left',
          paddingLeft: `${paddingLeft}px`,
          paddingRight: '8px',
          py: '12px',
          textTransform: 'none',
          width: '100%',
          ...(active && {
            color: theme.palette.secondary.main,
            fontWeight: 'bold',
            '& svg': {
              color: theme.palette.primary.main
            }
          })
        }}
        variant="text"
        to={path}
      >
        <Box style={{ flexGrow: 1 }}>
          {title}
        </Box>
        {info}
      </Button>
    </ListItem>
  );
};

NavItem.propTypes = {
  active: PropTypes.bool,
  children: PropTypes.node,
  depth: PropTypes.number.isRequired,
  icon: PropTypes.node,
  info: PropTypes.node,
  open: PropTypes.bool,
  path: PropTypes.string,
  title: PropTypes.string.isRequired
};

NavItem.defaultProps = {
  active: false,
  open: false
};

export default NavItem;

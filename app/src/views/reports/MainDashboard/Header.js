import React, {
  useRef,
  useState
} from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Breadcrumbs,
  Button,
  Grid,
  Link,
  Menu,
  MenuItem,
  SvgIcon,
  Typography,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import { Calendar as CalendarIcon } from 'react-feather';

const timeRanges = [
  {
    value: 'today',
    text: 'Today'
  },
  {
    value: 'yesterday',
    text: 'Yesterday'
  },
  {
    value: 'last_30_days',
    text: 'Last 30 days'
  },
  {
    value: 'last_year',
    text: 'Last year'
  }
];

const useStyles = makeStyles(() => ({
  root: {}
}));

const Header = ({ className, ...rest }) => {
  const classes = useStyles();
  const actionRef = useRef(null);
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [timeRange, setTimeRange] = useState(timeRanges[2].text);

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
          <Typography
            variant="body1"
            color="textPrimary"
          >
            Overivew
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          Here&apos;s what&apos;s happening
        </Typography>
      </Grid>
      <Grid item>
        <Typography
          variant="h5"
          color="textPrimary"
        >
          LAST 15 MINUTES
        </Typography>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string
};

export default Header;

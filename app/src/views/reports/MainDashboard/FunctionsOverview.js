import React from 'react';
import PropTypes from 'prop-types';
import { Link as RouterLink } from 'react-router-dom';
import Skeleton from '@material-ui/lab/Skeleton';
import clsx from 'clsx';
import {
  Card,
  Grid,
  Link,
  Typography,
  makeStyles
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {},
  item: {
    padding: theme.spacing(3),
    textAlign: 'center',
    [theme.breakpoints.up('md')]: {
      '&:not(:last-of-type)': {
        borderRight: `1px solid ${theme.palette.divider}`
      }
    },
    [theme.breakpoints.down('sm')]: {
      '&:not(:last-of-type)': {
        borderBottom: `1px solid ${theme.palette.divider}`
      }
    }
  },
  valueContainer: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  label: {
    marginLeft: theme.spacing(1)
  }
}));

const Overview = ({
  className,
  totalRequests,
  totalRequestsSinceEpoch,
  totalFunctions,
  totalImages,
  dbStats,
  ...rest
}) => {
  const classes = useStyles();

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Grid
        alignItems="center"
        container
        justify="space-between"
      >
        <Grid
          className={classes.item}
          item
          md={3}
          sm={6}
          xs={12}
        >
          {totalRequestsSinceEpoch === null ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Requests Made Since Epoch
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {totalRequestsSinceEpoch}
                </Typography>
              </div>
            </>
          }
        </Grid>
        <Grid
          className={classes.item}
          item
          md={3}
          sm={6}
          xs={12}
        >
          {totalRequests === null ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Requests Made
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {totalRequests}
                </Typography>
              </div>
            </>
          }
        </Grid>
        <Grid
          className={classes.item}
          item
          md={3}
          sm={6}
          xs={12}
        >
          {totalImages === null ? <Skeleton variant="rect" />
            :
            <>
              <Link
                underline={"always"}
                style={{ margin: "auto" }}
                component={RouterLink}
                to={`/app/images`}
              >
                <Typography
                  component="h2"
                  gutterBottom
                  variant="overline"
                  color="primary"
                >
                  Image Count
                </Typography>
              </Link>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {totalImages}
                </Typography>
              </div>
            </>
          }
        </Grid>
        <Grid
          className={classes.item}
          item
          md={3}
          sm={6}
          xs={12}
        >
          {totalFunctions === null ? <Skeleton variant="rect" />
            :
            <>
              <Link
                underline={"always"}
                style={{ margin: "auto" }}
                component={RouterLink}
                to={`/app/images`}
              >
                <Typography
                  component="h2"
                  gutterBottom
                  variant="overline"
                  color="primary"
                >
                  Function Count
              </Typography>
              </Link>
              <div className={classes.valueContainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {totalFunctions}
                </Typography>
              </div>
            </>
          }
        </Grid>
      </Grid>
    </Card >
  );
};

Overview.propTypes = {
  className: PropTypes.string
};

export default Overview;

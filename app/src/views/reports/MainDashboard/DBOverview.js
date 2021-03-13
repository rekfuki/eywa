import React from 'react';
import PropTypes from 'prop-types';
import { Link as RouterLink } from 'react-router-dom';
import clsx from 'clsx';
import {
  Card,
  Grid,
  Typography,
  makeStyles
} from '@material-ui/core';
import Skeleton from '@material-ui/lab/Skeleton';

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

function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

const DBOverview = ({
  className,
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
          {!dbStats ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Database Collection Count
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {dbStats.collection_count}
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
          {!dbStats ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Database Data Size
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {formatBytes(dbStats.data_size)}
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
          {!dbStats ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Database Average Object Size
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {formatBytes(dbStats.data_size)}
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
          {!dbStats ? <Skeleton variant="rect" />
            :
            <>
              <Typography
                component="h2"
                gutterBottom
                variant="overline"
                color="textSecondary"
              >
                Database Object Count
              </Typography>
              <div className={classes.valuecontainer}>
                <Typography
                  variant="h3"
                  color="textPrimary"
                >
                  {dbStats.objects}
                </Typography>
              </div>
            </>
          }
        </Grid>
      </Grid>
    </Card >
  );
};

export default DBOverview;

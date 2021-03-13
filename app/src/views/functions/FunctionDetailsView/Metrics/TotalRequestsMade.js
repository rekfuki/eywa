import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { useSnackbar } from 'notistack';
import {
  Card,
  Grid,
  Typography,
  makeStyles
} from '@material-ui/core';
import axios from 'src/utils/axios';

const useStyles = makeStyles((theme) => ({
  root: {
    minHeight: 300,
    height: "100%"
  },
  item: {
    margin: "auto",
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
  label: {
    marginLeft: theme.spacing(1)
  },
  overline: {
    marginTop: theme.spacing(1)
  }
}));

const TotalRequests = ({ functionId, endTime }) => {
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const [totalRequests, setTotalRequests] = useState(0);

  const getSeries = async () => {
    try {
      const payload = {
        "type": "instant",
        "series": ["gateway_function_invocation_total"],
        "group_by": "function_id",
        "label_matchers": `function_id="${functionId}"`,
        "query": `sum(<<index .Series 0>>{<<.LabelMatchers>>}) by(<<.GroupBy>>)`
      }

      const response = await axios.post(`/eywa/api/metrics/query`, payload);
      const data = response.data.Data;
      setTotalRequests(data.result.length > 0 ? data.result[0].value[1] : 0);

    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get metrics', {
        variant: 'error'
      });
    }
  };

  useEffect(() => {
    getSeries();
  }, [endTime])

  return (
    <Card className={classes.root}>
      <Grid
        alignItems="center"
        style={{ height: "100%" }}
        container
      >
        <Grid
          className={classes.item}
          item
          xs={12}

        >
          <Typography
            className={classes.overline}
            variant="overline"
            color="textSecondary"
          >
            Total Requests
          </Typography>
          <Typography
            variant="h2"
            color="textPrimary"
          >
            {totalRequests}
          </Typography>
        </Grid>
      </Grid>
    </Card>
  );
};

TotalRequests.prototype = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
}

export default TotalRequests;
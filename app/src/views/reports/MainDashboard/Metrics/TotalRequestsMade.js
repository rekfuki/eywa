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

const TotalRequests = ({ endTime }) => {
  const { enqueueSnackbar } = useSnackbar();
  const [totalRequests, setTotalRequests] = useState(0);

  const getSeries = async () => {
    try {
      const payload = {
        "type": "instant",
        "series": ["gateway_function_invocation_total"],
        "group_by": "user_id",
        "query": `sum(<<index .Series 0>>) by(<<.GroupBy>>)`
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
    <>
      <Typography
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
    </>
  );
};

TotalRequests.prototype = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
}

export default TotalRequests;
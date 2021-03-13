import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import {
  Card,
  CardContent,
  Grid
} from '@material-ui/core';
import RequestAndErrorRates from './RequestAndErrorRates';
import RequestRateByMethod from './RequestRateByMethod';
import RequestDurationPercentage from './RequestDurationPercentage';
import Top3APICallsByPath from './TopAPICallsByPath';
import TotalRequests from './TotalRequestsMade';
import ErrorRate from './ErrorRate';

const Metrics = ({
  functionId
}) => {
  const [time, setTime] = useState(null);

  useEffect(() => {
    let interval = setInterval(() => setTime(Date.now()), refreshInterval);
    return () => { if (interval) clearInterval(interval) }
  }, [])

  return (
    <Card>
      <CardContent >
        <Grid item container xs={12}>
          <Grid item xs={12} md={6}>
            <TotalRequests functionId={functionId} endTime={endTime} />
          </Grid>
          <Grid item xs={12} md={6}>
            <ErrorRate functionId={functionId} endTime={endTime} range={range} width={width} />
          </Grid>
          <Grid item xs={12} xl={6}>
            <RequestAndErrorRates functionId={functionId} endTime={endTime} range={range} width={width} />
          </Grid>
          <Grid item xs={12} xl={6}>
            <RequestRateByMethod functionId={functionId} endTime={endTime} range={range} width={width} />
          </Grid>
          <Grid item xs={12} xl={6}>
            <RequestDurationPercentage functionId={functionId} endTime={endTime} range={range} width={width} />
          </Grid>
          <Grid item xs={12} xl={6}>
            <Top3APICallsByPath functionId={functionId} endTime={endTime} range={range} width={width} />
          </Grid>
        </Grid >
      </CardContent>
    </Card >
  );
};

Metrics.propTypes = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
};

export default Metrics;

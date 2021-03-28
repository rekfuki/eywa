import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { useSnackbar } from 'notistack';
import useDimensions from "react-use-dimensions";
import clsx from 'clsx';
import moment from 'moment';
import {
  Box,
  Card,
  CardContent,
  Divider,
  Grid,
  FormControl,
  MenuItem,
  IconButton,
  InputAdornment,
  InputLabel,
  TextField,
  Tooltip,
  Typography,
  Select,
  useTheme,
  makeStyles
} from '@material-ui/core';
import {
  Plus as PlusIcon,
  Minus as MinusIcon,
  ChevronRight as ChevronRightIcon,
  ChevronLeft as ChevronLeftIcon,
  X as XIcon
} from 'react-feather';
import {
  formatTime,
  formatTimeNoSeconds,
  formatDuration,
  roundInterval
} from 'src/utils/time';
import RequestAndErrorRates from './RequestAndErrorRates';
import RequestRateByMethod from './RequestRateByMethod';
import RequestDurationPercentage from './RequestDurationPercentage';
import Top3APICallsByPath from './Top3APICallsByPath';
import TotalRequests from './TotalRequestsMade';
import ErrorRate from './ErrorRate';
import QueueDwellHeatmap from './QueueDwellHeatmap';

const useStyles = makeStyles((theme) => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap'
  },
  textField: {
    maxWidth: "245px"
  },
  underline: {
    "&&&:before": {
      borderBottomStyle: "solid"
    }
  }
}));


const Metrics = ({
  functionId,
  className,
  ...rest
}) => {
  const classes = useStyles();

  const { enqueueSnackbar } = useSnackbar();

  const [range, setRange] = useState(5 * 60 * 1000)
  const [endTime, setEndTime] = useState(null);
  const [refreshInterval, setRefreshInterval] = useState(5000);
  const [time, setTime] = useState(null);

  const [measureRef, { width }] = useDimensions();

  let interval;

  const rangeSteps = [
    5 * 60,
    15 * 60,
    30 * 60,
    60 * 60,
    2 * 60 * 60,
    6 * 60 * 60,
    12 * 60 * 60,
    24 * 60 * 60,
    48 * 60 * 60,
    7 * 24 * 60 * 60,
    14 * 24 * 60 * 60,
    28 * 24 * 60 * 60,
    56 * 24 * 60 * 60,
    365 * 24 * 60 * 60,
    730 * 24 * 60 * 60
  ].map(s => s * 1000);

  const increaseRange = () => {
    for (const rv of rangeSteps) {
      if (range < rv) {
        setRange(rv);
        return;
      }
    }
  };

  const decreaseRange = () => {
    for (const rv of rangeSteps.slice().reverse()) {
      if (range > rv) {
        setRange(rv);
        return;
      }
    }
  };

  const getEndTime = () => {
    return endTime || moment().valueOf();
  };

  const clearEndTime = () => {
    setEndTime(null);
  }

  const getRange = () => {
    return range;
  }

  const getWidth = () => {
    return width;
  }

  const getRefreshInterval = () => {
    return refreshInterval;
  }

  const increaseTime = () => {
    const updatedTime = getEndTime() + getRange();
    setEndTime(updatedTime);
    console.log(formatTime(updatedTime))
  };

  const decreaseTime = () => {
    const updatedTime = getEndTime() - getRange();
    setEndTime(updatedTime);
    console.log(formatTime(updatedTime))
  };

  useEffect(() => {
    if (interval) {
      clearInterval(interval);
      interval = null;
    }

    interval = setInterval(() => setTime(Date.now()), refreshInterval);
    return () => { if (interval) clearInterval(interval) }
  }, [refreshInterval])

  return (
    <Card>
      <CardContent >
        <Grid
          container
          spacing={3}
          ref={measureRef}
        >
          <Grid item xs={12} sm={"auto"}>
            <Box>
              <FormControl style={{ minWidth: "110px" }}>
                <InputLabel id="select-label">Time Range</InputLabel>
                <Select
                  labelId="select-label"
                  value={range}
                  onChange={(event) => setRange(event.target.value)}
                >
                  <MenuItem value={1 * 60 * 1000}>1m</MenuItem>
                  <MenuItem value={5 * 60 * 1000}>5m</MenuItem>
                  <MenuItem value={15 * 60 * 1000}>15m</MenuItem>
                  <MenuItem value={30 * 60 * 1000}>30m</MenuItem>
                  <MenuItem value={60 * 60 * 1000}>1h</MenuItem>
                  <MenuItem value={2 * 60 * 60 * 1000}>2h</MenuItem>
                  <MenuItem value={6 * 60 * 60 * 1000}>6h</MenuItem>
                  <MenuItem value={12 * 60 * 60 * 1000}>12h</MenuItem>
                  <MenuItem value={24 * 60 * 60 * 1000}>24h</MenuItem>
                  <MenuItem value={48 * 60 * 60 * 1000}>48h</MenuItem>
                  <MenuItem value={7 * 24 * 60 * 60 * 1000}>7d</MenuItem>
                </Select>
              </FormControl>
              {/* <Tooltip title="Decrease range">
                <IconButton onClick={decreaseRange} aria-label="decrease">
                  <MinusIcon />
                </IconButton>
              </Tooltip>
              <TextField
                label="Range"
                style={{ maxWidth: "70px" }}
                id="outlined-read-only-input"
                value={formatDuration(range)}
                disabled
                inputProps={{
                  style: { color: 'black' }
                }}
                InputProps={{
                  readOnly: true,
                  className: classes.underline
                }}
              />
              <Tooltip title="Increase range">
                <IconButton onClick={increaseRange} aria-label="increase">
                  <PlusIcon />
                </IconButton>
              </Tooltip> */}
            </Box>
          </Grid>
          <Grid item style={{ display: "flex" }}>
            <Box display="flex" xs={12} sm={"auto"}>
              <Tooltip title="Decrease time">
                <IconButton onClick={decreaseTime} aria-label="decrease">
                  <ChevronLeftIcon />
                </IconButton>
              </Tooltip>
              <form className={classes.container} noValidate>
                <TextField
                  label="Time"
                  id="datetime-local"
                  type={endTime ? "dateendTime-local" : "text"}
                  value={endTime ? formatTimeNoSeconds(endTime) : "End Time"}
                  placeholder="Now"
                  className={classes.textField}
                  onClick={!endTime ? () => setEndTime(getEndTime()) : undefined}
                  InputProps={{
                    readOnly: endTime ? false : true,
                    endAdornment: !endTime ? undefined : (
                      <InputAdornment position="end">
                        <IconButton onClick={clearEndTime} size="small" style={{ borderRadius: 0 }} >
                          <XIcon style={{ paddingBottom: "2px" }} />
                        </IconButton>
                      </InputAdornment>
                    )
                  }}
                  InputLabelProps={{
                    shrink: true
                  }}
                />
              </form>
              <Tooltip title="Increase time">
                <IconButton onClick={increaseTime} aria-label="increase">
                  <ChevronRightIcon />
                </IconButton>
              </Tooltip>
            </Box>
          </Grid>
          <Grid item style={{ textAlign: 'center' }}>
            <Box flexGrow={1} alignSelf="flex-end" justifyContent="flex-end">
              <FormControl style={{ minWidth: "110px" }}>
                <InputLabel id="select-label">Update interval</InputLabel>
                <Select
                  labelId="select-label"
                  value={refreshInterval}
                  onChange={(event) => setRefreshInterval(event.target.value)}
                >
                  <MenuItem value={5000}>5s</MenuItem>
                  <MenuItem value={10000}>10s</MenuItem>
                  <MenuItem value={30000}>30s</MenuItem>
                  <MenuItem value={60000}>1m</MenuItem>
                  <MenuItem value={300000}>5m</MenuItem>
                  <MenuItem value={600000}>10m</MenuItem>
                </Select>
              </FormControl>
            </Box>
          </Grid>
          <Grid item>
            <Divider></Divider>
          </Grid>
        </Grid>
        {(() => {
          const endTime = getEndTime();
          const range = getRange();
          const refreshInterval = getRefreshInterval();
          const width = getWidth();

          if (!width) {
            return null
          }

          return (
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
          );
        })()}

      </CardContent>
    </Card >
  );
};

Metrics.propTypes = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
};

export default Metrics;

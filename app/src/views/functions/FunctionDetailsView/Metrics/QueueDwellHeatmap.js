import React, { Fragment, useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import useDimensions from "react-use-dimensions";

import Chart from 'react-apexcharts';
import moment from 'moment';
import { useSnackbar } from 'notistack';
import {
  Box,
  Card,
  CardContent,
  Grid,
  LinearProgress,
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
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import axios from 'src/utils/axios';
import {
  formatTime,
  formatTimeNoSeconds,
  formatDuration,
  roundInterval
} from 'src/utils/time';
import { parseValue } from 'src/utils/numbers';

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


const QueueDwellHeatMap = ({ functionId }) => {
  const theme = useTheme();
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();

  const { enqueueSnackbar } = useSnackbar();

  const [range, setRange] = useState(5 * 60 * 1000)
  const [endTime, setEndTime] = useState(null);
  const [refreshInterval, setRefreshInterval] = useState(5000);
  const [chart, setChart] = useState(null);
  const [loading, setLoading] = useState(true);

  const [measureRef, { width }] = useDimensions();

  const minInterval = 5 * 1000;

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

  const getSeries = async () => {
    try {
      const range = getRange()
      const width = getWidth();
      let step = roundInterval(range / 25);
      step = step < minInterval ? minInterval : step;
      const endTime = Math.floor(getEndTime() / step) * step;
      const startTime = endTime - range;

      const payload = {
        "type": "range",
        "series": ["gateway_queue_dwell_duration_milliseconds_bucket"],
        "group_by": "le",
        "label_matchers": `function_id="${functionId}"`,
        "query": `sum(rate(<<index .Series 0>>{<<.LabelMatchers>>}[${step}ms])) by(<<.GroupBy>>)`,
        "start": startTime / 1000,
        "end": endTime / 1000,
        "step": step / 1000
      }

      const response = await axios.post(`/eywa/api/metrics/query`, payload);
      const series = processMetrics(response.data.Data);

      const newChart = buildChart(series);
      addTimeAxis(newChart.options, getWidth(), endTime);

      console.log(newChart);
      setChart(newChart);

    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get metrics', {
        variant: 'error'
      });
    }
    setLoading(false);
  };

  const processMetrics = (data, label = false) => {
    const series = data.result.map(({ values, metric }) => {
      const data = []
      for (let i = 0; i < values.length; i++) {
        data.push([values[i][0] * 1000, parseValue(values[i][1])])
      }
      return {
        label: metric.le,
        name: metric.le + " ms",
        data: data
      }
    });
    return series.sort((sortSeriesByLabel));
  }

  function sortSeriesByLabel(s1, s2) {
    let label1, label2;

    try {
      // fail if not integer. might happen with bad queries
      label1 = parseHistogramLabel(s1.label);
      label2 = parseHistogramLabel(s2.label);
    } catch (err) {
      console.error(err.message || err);
      return 0;
    }

    if (label1 > label2) {
      return 1;
    }

    if (label1 < label2) {
      return -1;
    }

    return 0;
  }

  function parseHistogramLabel(label) {
    if (label === '+Inf' || label === 'inf') {
      return +Infinity;
    }
    const value = Number(label);
    if (isNaN(value)) {
      throw new Error(`Error parsing histogram label: ${label} is not a number`);
    }
    return value;
  }

  useEffect(() => {
    if (width) {
      if (interval) {
        clearInterval(interval);
        interval = null;
      }

      getSeries();

      interval = setInterval(getSeries, refreshInterval);
    }

    return () => { if (interval) clearInterval(interval) }
  }, [range, endTime, refreshInterval, width])


  const graphTimeFormat = (ticks, min, max) => {
    if (min && max && ticks) {
      const range = max - min;
      const secPerTick = range / ticks / 1000;
      // Need have 10 millisecond margin on the day range
      // As sometimes last 24 hour dashboard evaluates to more than 86400000
      const oneDay = 86400010;
      const oneYear = 31536000000;

      if (secPerTick <= 45) {
        return 'HH:mm:ss';
      }
      if (secPerTick <= 7200 || range <= oneDay) {
        return 'HH:mm';
      }
      if (secPerTick <= 80000) {
        return 'MM/DD HH:mm';
      }
      if (secPerTick <= 2419200 || range <= oneYear) {
        return 'MM/DD';
      }
      if (secPerTick <= 31536000) {
        return 'YYYY-MM';
      }
      return 'YYYY';
    }

    return 'HH:mm';
  };

  const addTimeAxis = (options, width, endTime) => {
    const ticks = width ? width / 100 : 2;

    const min = endTime - getRange();
    const max = endTime;


    options.xaxis = {
      min: min,
      max: max,
      label: 'Datetime',
      labels: {
        formatter: (_, timestamp) => {
          return moment(timestamp).format(graphTimeFormat(ticks, min, max))
        },
        style: {
          colors: theme.palette.text.secondary
        }
      },
      tickAmount: ticks,
      tickPlacement: 'on',
      axisBorder: {
        color: theme.palette.divider
      },
      axisTicks: {
        show: true,
        color: theme.palette.divider
      }
    };
  }

  const buildChart = (series) => {
    return {
      series: series,
      type: 'heatmap',
      options: {
        plotOptions: {
          heatmap: {
            shadeIntensity: 0.5,

            colorScale: {
              ranges: [{
                from: 0,
                to: 0,
                color: theme.palette.text.secondary
              }
              ]
            }
          }
        },
        noData: {
          text: "No data available"
        },
        chart: {
          stacked: true,
          background: theme.palette.background.paper,
          toolbar: {
            show: false
          },
          animations: {
            enabled: false
          },
          zoom: {
            enabled: false
          }
        },
        dataLabels: {
          enabled: false
        },
        grid: {
          xaxis: {
            lines: {
              show: true
            }
          },
          yaxis: {
            lines: {
              show: true
            }
          },
          borderColor: theme.palette.divider
        },
        legend: {
          show: false
        },
        markers: {
          size: 0
        },
        stroke: {
          width: 1,
          curve: 'straight',
          lineCap: 'butt'
        },
        title: {
          text: "Queue Dwell Time (async)",
          align: "center"
        },
        theme: {
          mode: theme.palette.type
        },
        tooltip: {
          theme: theme.palette.type,
          x: {
            formatter: (value) => (moment(value).format('dd/MM/yy HH:mm'))
          }
        },
        xaxis: [],
        yaxis: {
          decimalsInFloat: 2,
          axisTicks: {
            show: true,
            color: theme.palette.divider
          },
          axisBorder: {
            show: true,
            color: theme.palette.divider
          },
          labels: {
            style: {
              colors: theme.palette.text.secondary
            }
          }
        }
      }
    };
  }

  return (
    <Card>
      <CardContent >
        <Typography
          variant="h4"
          color="textPrimary"
        >
          Requests Per Second
        </Typography>
        <Box flexGrow={1} mt={3} />
        <Grid
          container
          spacing={3}
          ref={measureRef}
        >

          <Grid item xs={12} sm={"auto"}>
            <Box>
              <Tooltip title="Decrease range">
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
              </Tooltip>
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
        </Grid>
        {loading ? <LinearProgress /> :
          <Fragment>
            <Grid item >
              {chart && <Chart
                type="line"
                height="300"
                {...chart}
              />}
            </Grid>
          </Fragment>
        }
      </CardContent>
    </Card >
  );
};

QueueDwellHeatMap.prototype = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
}

export default QueueDwellHeatMap;
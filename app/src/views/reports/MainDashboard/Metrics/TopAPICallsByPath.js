import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import Chart from 'react-apexcharts';
import moment from 'moment';
import {
  Card,
  CardContent,
  LinearProgress,
  useTheme
} from '@material-ui/core';

const TopAPICallsByPath = ({ series, width, endTime, startTime }) => {
  const theme = useTheme();
  const [chart, setChart] = useState(null);

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

  const addTimeAxis = (options, width, endTime, startTime) => {
    const ticks = width ? width / 100 : 2;

    const min = startTime;
    const max = endTime;

    options.xaxis = {
      type: "datetime",
      min: min,
      max: max,
      label: 'Datetime',
      labels: {
        formatter: (value, timestamp) => {
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
      type: 'area',
      options: {
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
          show: true,
          showForSingleSeries: true,
          position: 'bottom',
          horizontalAlign: 'right',
          labels: {
            colors: theme.palette.text.secondary
          }
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
          text: "Top 5 API calls (by path)",
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

  useEffect(() => {
    if (!series || !width || !endTime || !startTime) {
      return null
    }
    const chart = buildChart(series);
    addTimeAxis(chart.options, width, endTime, startTime);
    setChart(chart);
  }, [series, width, endTime, startTime])

  return (
    <Card>
      <CardContent >
        {!chart
          ? <LinearProgress />
          : <Chart
            type="line"
            height="300"
            {...chart}
          />}
      </CardContent>
    </Card >
  );
};

TopAPICallsByPath.prototype = {
  className: PropTypes.string
}

export default TopAPICallsByPath;
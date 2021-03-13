import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import Chart from 'react-apexcharts';
import { useSnackbar } from 'notistack';
import {
  Card,
  CardContent,
  useTheme
} from '@material-ui/core';
import axios from 'src/utils/axios';
import {
  roundInterval
} from 'src/utils/time';

const ErrorRate = ({ functionId, endTime, range, width }) => {
  const theme = useTheme();
  const { enqueueSnackbar } = useSnackbar();
  const [series, setSeries] = useState([]);
  const minInterval = 5 * 1000;


  const getSeries = async () => {
    try {
      let step = roundInterval(range / width);
      step = step < minInterval ? minInterval : step;
      endTime = Math.floor(endTime / step) * step;

      const payload = {
        "type": "instant",
        "series": ["gateway_function_invocation_total"],
        "group_by": "function_id",
        "label_matchers": `function_id="${functionId}",code=~"4..|5.."`,
        "query": `((sum by(<<.GroupBy>>) (rate(<<index .Series 0>>{<<.LabelMatchers>>}[${step}ms]))) / (sum by(<<.GroupBy>>) (rate(<<index .Series 0>>[${step}ms])))) * 100`,
        "end": endTime / 1000
      }

      const response = await axios.post(`/eywa/api/metrics/query`, payload);
      const data = response.data.Data;
      setSeries([data.result.length == 0 ? [0] : isNaN(data.result[0].value[1]) ? [0] : data.result[0].value[1]]);

    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get metrics', {
        variant: 'error'
      });
    }
  };

  useEffect(() => {
    getSeries();
  }, [endTime, range, width])

  if (!series) {
    return null
  }

  const chart = {
    series: series,
    type: 'radialBar',
    options: {
      noData: {
        text: "No data available"
      },
      chart: {
        background: theme.palette.background.paper,
        toolbar: {
          show: false
        },
        animations: {
          enabled: false
        }
      },
      title: {
        text: "Errors %",
        align: "center"
      },
      theme: {
        mode: theme.palette.type
      },
      legend: {
        show: false
      },
      plotOptions: {
        radialBar: {
          startAngle: -135,
          endAngle: 135,
          dataLabels: {
            name: {
              show: false
            },
            value: {
              offsetY: 76,
              fontSize: '22px',
              color: undefined,
              formatter: function (val) {
                return Math.round(val) + "%";
              }
            }
          }
        }
      },
      fill: {
        type: 'gradient',
        gradient: {
          shade: 'dark',
          shadeIntensity: 0.15,
          inverseColors: false,
          opacityFrom: 1,
          opacityTo: 1,
          stops: [0, 50, 65, 91]
        }
      },
      colors: [theme.palette.error.main],
      stroke: {
        dashArray: 4
      },
      states: {
        hover: {
          filter: {
            type: 'none'
          }
        }
      }
    }
  }

  return (
    <Card>
      <CardContent >
        {chart && <Chart
          type="line"
          height="300"
          {...chart}
        />}
      </CardContent>
    </Card >
  );
};

ErrorRate.prototype = {
  className: PropTypes.string,
  functionId: PropTypes.string.isRequired
}

export default ErrorRate;
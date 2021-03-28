import React, { useState, useEffect } from 'react';
import useDimensions from "react-use-dimensions";
import moment from 'moment';
import { useSnackbar } from 'notistack';
import {
  Container,
  Grid,
  makeStyles
} from '@material-ui/core';
import {
  roundInterval
} from 'src/utils/time';
import { parseValue } from 'src/utils/numbers';
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import Header from './Header';
import FunctionOverview from './FunctionsOverview';
import DBOverview from './DBOverview';
import RequestAndErrorRates from './Metrics/RequestAndErrorRates';
import TopAPICallsByPath from './Metrics/TopAPICallsByPath';
import Logs from './Logs';
import Timelines from './Timelines';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const DashboardAlternativeView = () => {
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const [measureRef, { width }] = useDimensions();
  const [totalRequests, setTotalRequests] = useState(null);
  const [totalRequestsSinceEpoch, setTotalRequestsSinceEpoch] = useState(null);
  const [totalFunctions, setTotalFunctions] = useState(null);
  const [totalImages, setTotalImages] = useState(null);
  const [requestAndErrorRatesData, setRequestAndErrorRatesData] = useState({});
  const [top5Data, setTop5Data] = useState({});
  const [dbStats, setDBStats] = useState(null);
  const refreshInterval = 5000;
  const minInterval = 5000;
  const range = 15 * 60 * 1000;


  const getData = async () => {
    try {
      let step = roundInterval(range / width);
      step = step < minInterval ? minInterval : step;
      const endTime = Math.floor(moment().valueOf() / step) * step;
      const startTime = endTime - range;

      getTotalRequests();
      getTotalRequestsSinceEpoch();
      getFunctions();
      getImages();
      getDBStats();
      getTop5Paths(startTime, endTime, step);

      Promise.all([
        getAllRequests(startTime, endTime, step),
        getXXXRequests(`code=~"2.."`, startTime, endTime, step),
        getXXXRequests(`code=~"4.."`, startTime, endTime, step),
        getXXXRequests(`code=~"5.."`, startTime, endTime, step)
      ])
        .then(([allRequests, status2XX, status4XX, status5XX]) => {
          if (allRequests.length > 0) allRequests[0] = { ...allRequests[0], name: "Request Rate" }
          if (status2XX.length > 0) status2XX[0] = { ...status2XX[0], name: "2XX Success Rate" }
          if (status4XX.length > 0) status4XX[0] = { ...status4XX[0], name: "4XX Errors Rate" }
          if (status5XX.length > 0) status5XX[0] = { ...status5XX[0], name: "5XX Errors Rate" }
          const series = allRequests.concat(status2XX).concat(status4XX).concat(status5XX)

          setRequestAndErrorRatesData({
            series: series,
            width: width,
            endTime: endTime,
            startTime: startTime
          })
        })
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get metrics', {
        variant: 'error'
      });
    }
  };

  const getTop5Paths = async (startTime, endTime, step) => {
    const payload = {
      "type": "range",
      "series": ["gateway_function_invocation_total"],
      "group_by": "path",
      "query": `topk(3, sum(rate(<<index .Series 0>>[${step}ms])) by(<<.GroupBy>>))`,
      "start": startTime / 1000,
      "end": endTime / 1000,
      "step": step / 1000
    }
    const response = await axios.post(`/eywa/api/metrics/query`, payload);
    const series = processMetrics(response.data.Data, true, true);

    setTop5Data({
      series: series,
      width: width,
      endTime: endTime,
      startTime: startTime
    })
  };

  const getAllRequests = async (startTime, endTime, step) => {
    const payload = {
      "type": "range",
      "series": ["gateway_function_invocation_total"],
      "group_by": "user_id",
      "query": `sum(rate(<<index .Series 0>>[${step}ms])) by(<<.GroupBy>>)`,
      "start": startTime / 1000,
      "end": endTime / 1000,
      "step": step / 1000
    }
    const response = await axios.post(`/eywa/api/metrics/query`, payload)
    return processMetrics(response.data.Data)
  };

  const getXXXRequests = async (codeLabel, startTime, endTime, step) => {
    const payload = {
      "type": "range",
      "series": ["gateway_function_invocation_total"],
      "label_matchers": `${codeLabel}`,
      "group_by": "user_id",
      "query": `sum(rate(<<index .Series 0>>{<<.LabelMatchers>>}[${step}ms])) by(<<.GroupBy>>)`,
      "start": startTime / 1000,
      "end": endTime / 1000,
      "step": step / 1000
    }
    const response = await axios.post(`/eywa/api/metrics/query`, payload)
    return processMetrics(response.data.Data)
  }

  const processMetrics = (data, label = false) => {
    return data.result.map(({ values, metric }) => {
      const data = []
      for (let i = 0; i < values.length; i++) {
        data.push([values[i][0] * 1000, parseValue(values[i][1])])
      }
      return {
        name: label ? JSON.stringify(metric) : undefined,
        data: data
      }
    });
  }

  const getTotalRequests = async () => {
    const payload = {
      "type": "instant",
      "series": ["gateway_function_invocation_total"],
      "group_by": "user_id",
      "query": `sum by (<<.GroupBy>>) (round(increase(<<index .Series 0>>[15m])))`
    }
    const response = await axios.post(`/eywa/api/metrics/query`, payload);
    const data = response.data.Data;
    setTotalRequests(data.result.length > 0 ? data.result[0].value[1] : 0);
  };

  const getTotalRequestsSinceEpoch = async () => {
    const payload = {
      "type": "instant",
      "series": ["gateway_function_invocation_total"],
      "group_by": "user_id",
      "query": `sum by (<<.GroupBy>>) (<<index .Series 0>>)`
    }
    const response = await axios.post(`/eywa/api/metrics/query`, payload);
    const data = response.data.Data;
    setTotalRequestsSinceEpoch(data.result.length > 0 ? data.result[0].value[1] : 0);
  };

  const getFunctions = async () => {
    const response = await axios.get(`/eywa/api/functions`);
    setTotalFunctions(response.data.total_count)
  }

  const getImages = async () => {
    const response = await axios.get(`/eywa/api/images`);
    setTotalImages(response.data.total_count)
  }

  const getDBStats = async () => {
    const response = await axios.get(`/eywa/api/database`);
    setDBStats(response.data)
  }

  useEffect(() => {
    if (!width) {
      return
    }

    if (interval) {
      clearInterval(interval);
      interval = null;
    }

    getData();

    let interval = setInterval(getData, refreshInterval);
    return () => { if (interval) clearInterval(interval) }
  }, [width])

  return (
    <Page
      className={classes.root}
      title="Dashboard"
    >
      <Container maxWidth={false} ref={measureRef}>
        <Header />
        <Grid
          container
          spacing={3}
        >
          <Grid
            item
            xs={12}
          >
            <FunctionOverview
              totalRequests={totalRequests}
              totalRequestsSinceEpoch={totalRequestsSinceEpoch}
              totalFunctions={totalFunctions}
              totalImages={totalImages}
            />
            <DBOverview
              dbStats={dbStats}
            />
          </Grid>
          <Grid
            item
            xl={6}
            xs={12}
          >
            <RequestAndErrorRates {...requestAndErrorRatesData} />
          </Grid>
          <Grid
            item
            xl={6}
            xs={12}
          >
            <TopAPICallsByPath {...top5Data} />
          </Grid>
          <Grid
            item
            xs={12}
          >
            <Logs />
          </Grid>
          <Grid
            item
            xs={12}
          >
            <Timelines />
          </Grid>
        </Grid>
      </Container>
    </Page>
  );
};

export default DashboardAlternativeView;

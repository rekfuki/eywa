import React, {
  useCallback,
  useState,
  useEffect
} from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {
  Box,
  Container,
  Divider,
  Tab,
  Tabs,
  makeStyles
} from '@material-ui/core';
import Page from 'src/components/Page';
import axios from 'src/utils/axios';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Header from './Header/Header';
import Details from './Details';
import Logs from './Logs';
import Timelines from './Timelines';
import Metrics from './Metrics'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const FunctionDetailsView = () => {
  const { functionId } = useParams();
  const classes = useStyles();
  const history = useHistory();
  const isMountedRef = useIsMountedRef();
  const [fn, setFn] = useState(null);
  const [currentTab, setCurrentTab] = useState('details');
  const { enqueueSnackbar } = useSnackbar();
  const [timelines, setTimelines] = useState({});

  const tabs = [
    { value: 'details', label: 'Details' },
    { value: 'timelines', label: 'Timelines' },
    { value: 'logs', label: 'Logs' },
    { value: 'metrics', label: 'Metrics' }
  ];

  const onDeleteFn = async () => {
    try {
      await axios.delete('/eywa/api/functions/' + fn.id);
      enqueueSnackbar('Function deleted', {
        variant: 'success'
      });

      history.push("/app/functions")
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Delete request failed', {
        variant: 'error'
      });
    }
  }

  const handleTabsChange = (event, value) => {
    setCurrentTab(value);
  };

  const getTimelines = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/timeline?function_id=' + functionId)

      if (isMountedRef.current) {
        setTimelines(response.data);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get timelines', {
        variant: 'error'
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getTimelines();
  }, [getTimelines]);


  const getFn = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/functions/' + functionId);

      if (isMountedRef.current) {
        setFn(response.data);
      }
    } catch (err) {
      console.error(err);
      let message = "Failed to get function"
      if (err.response.status === 404) {
        message = "Function Not Found"
      }
      enqueueSnackbar(message, {
        variant: 'error'
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getFn();
  }, [getFn]);

  if (!fn || !timelines) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Function Details"
    >
      <Container maxWidth={false}>
        <Header fn={fn} onDeleteFn={onDeleteFn} />
        <Box mt={3}>
          <Tabs
            onChange={handleTabsChange}
            scrollButtons="auto"
            value={currentTab}
            variant="scrollable"
            textColor="secondary"
          >
            {tabs.map((tab) => (
              <Tab
                key={tab.value}
                label={tab.label}
                value={tab.value}
              />
            ))}
          </Tabs>
        </Box>
        <Divider />
        <Box mt={2}>
          {currentTab === 'details' && <Details fn={fn} />}
          {currentTab === 'timelines' && <Timelines timelines={timelines} functionId={functionId} />}
          {currentTab === 'logs' && <Logs functionId={functionId} />}
          {currentTab === 'metrics' && <Metrics functionId={functionId} />}
        </Box>
      </Container>
    </Page>
  );
};

export default FunctionDetailsView;

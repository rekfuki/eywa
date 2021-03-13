import React, {
  useCallback,
  useState,
  useEffect
} from 'react';
import { useParams } from 'react-router-dom';
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
import Header from './Header';
import Details from './Details';
import Logs from './Logs';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const TimelineDetailsView = () => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const { requestId } = useParams();
  const { enqueueSnackbar } = useSnackbar();
  const [timeline, setTimeline] = useState(null);
  const [currentTab, setCurrentTab] = useState('details');

  const tabs = [
    { value: 'details', label: 'Details' },
    { value: 'logs', label: 'Logs' }
  ];

  const getTimeline = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/timeline/'+requestId);

      if (isMountedRef.current) {
        setTimeline(response.data);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get timeline details', {
        variant: 'error',
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getTimeline();
  }, [getTimeline]);

  const handleTabsChange = (event, value) => {
    setCurrentTab(value);
  };

  if (!timeline) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Execution Details"
    >
      <Container maxWidth={false}>
        <Header requestId={requestId}/>
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
          {currentTab === 'details' && <Details timeline={timeline} />}
          {currentTab === 'logs' && <Logs requestId={requestId}/>}
        </Box>
      </Container>
    </Page>
  );
};

export default TimelineDetailsView;

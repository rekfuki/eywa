import React, {
  useState,
  useEffect,
  useCallback
} from 'react';
import {
  Box,
  Container,
  makeStyles
} from '@material-ui/core';
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Header from './Header';
import Results from './Results';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const FunctionsListView = () => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const [fns, setFns] = useState([]);

  const getFns = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/functions');
      if (isMountedRef.current) {
        setFns(response.data.objects);
      }
    } catch (err) {
      console.error(err);
    }
  }, [isMountedRef]);

  useEffect(() => {
    getFns();
  }, [getFns]);

  return (
    <Page
      className={classes.root}
      title="Functions"
    >
      <Container maxWidth={false}>
        <Header />
        <Box mt={3}>
          <Results fns={fns} />
        </Box>
      </Container>
    </Page>
  );
};

export default FunctionsListView;

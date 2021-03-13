import React, {
  useState,
  useCallback,
  useEffect
} from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Container,
  makeStyles
} from '@material-ui/core';
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import FnEditForm from './FunctionEditForm';
import Header from './Header';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const FunctionEditView = () => {
  const { functionId } = useParams();
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const [fn, setFn] = useState(null);
  const [secrets, setSecrets] = useState(null);
  const [images, setImages] = useState(null);

  const getFn = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/functions/' + functionId);
      if (isMountedRef.current) {
        setFn(response.data);
      }
    } catch (err) {
      console.error(err);
    }
  }, [isMountedRef]);

  const getSecrets = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/secrets');
      if (isMountedRef.current) {
        setSecrets(response.data.objects);
      }
    } catch (err) {
      console.error(err);
    }
  }, [isMountedRef]);

  const getImages = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/images?per_page=1000');
      if (isMountedRef.current) {
        setImages(response.data.objects);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get images', {
        variant: 'error'
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getImages();
  }, [getImages]);

  useEffect(() => {
    getFn();
  }, [getFn]);

  useEffect(() => {
    getSecrets();
  }, [getSecrets]);

  if (!fn || !secrets || !images) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Function Edit"
    >
      <Container maxWidth={false}>
        <Header fn={fn} />
      </Container>
      <Box mt={3}>
        <Container maxWidth="lg">
          <FnEditForm fn={fn} secrets={secrets} images={images} />
        </Container>
      </Box>
    </Page>
  );
};

export default FunctionEditView;

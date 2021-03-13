import React, {
  useState,
  useCallback,
  useEffect
} from 'react';
import { useParams } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {
  Box,
  Container,
  makeStyles
} from '@material-ui/core';
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import FnEditForm from './SecretEditForm';
import Header from './Header';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const SecretEditView = () => {
  const { secretId } = useParams();
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const isMountedRef = useIsMountedRef();
  const [secret, setSecret] = useState(null);

  const getSecret = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/secrets/' + secretId);
      if (isMountedRef.current) {
        setSecret(response.data);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get secrets', {
        variant: 'error'
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getSecret();
  }, [getSecret]);

  if (!secret) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Function Edit"
    >
      <Container maxWidth={false}>
        <Header secret={secret} />
      </Container>
      <Box mt={3}>
        <Container maxWidth="lg">
          <FnEditForm secret={secret} />
        </Container>
      </Box>
    </Page>
  );
};

export default SecretEditView;

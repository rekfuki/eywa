import React, {
  useCallback,
  useEffect,
  useState
} from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {
  Box,
  Container,
  Grid,
  makeStyles
} from '@material-ui/core';
import axios from 'src/utils/axios';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Page from 'src/components/Page';
import Header from './Header/Header';
import SecretInfo from './SecretInfo';
import Fields from './Fields';
import Mounts from './Mounts';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const OrderDetailsView = () => {
  const classes = useStyles();
  const history = useHistory();
  const { secretId } = useParams();
  const { enqueueSnackbar } = useSnackbar();
  const isMountedRef = useIsMountedRef();
  const [secret, setSecret] = useState(null);

  const onDeleteSecret = async () => {
    try {
      await axios.delete('/eywa/api/secrets/' + secret.id);
      enqueueSnackbar('Secret deleted', {
        variant: 'success'
      });

      history.push("/app/secrets")
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Delete request failed', {
        variant: 'error'
      });
    }
  }

  const getSecret = async () => {
    try {
      const response = await axios.get('/eywa/api/secrets/' + secretId);
      console.log("Retreived secret data: ", response.data);
      setSecret(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get secret', {
        variant: 'error'
      });
    }
  };

  const updateSecret = async (payload) => {
    try {
      const response = await axios.put("/eywa/api/secrets/" + secret.id, payload);

      if (isMountedRef.current) {
        // setSecret({ ...secret, ...response.data });
        await getSecret();
        enqueueSnackbar('Secret updated', {
          variant: 'success'
        });
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to update', {
        variant: 'error'
      });
    }
  }

  useEffect(() => {
    getSecret();
  }, []);


  if (!secret) {
    return null;
  }

  const disabledEditing = secret.name.includes("mongodb")
  return (
    <Page
      className={classes.root}
      title="Secret Details"
    >
      <Container maxWidth={false}>
        <Header secret={secret} disabledEditing={disabledEditing} onDeleteSecret={onDeleteSecret} />
        <Box mt={2}>
          <Grid
            container
            spacing={3}
          // direction="row"
          // // justify="center"
          // // alignItems="baseline"

          >
            <Grid
              item
              xl={3}
              md={6}
              xs={12}
            >
              <SecretInfo secret={secret} />
            </Grid>
            <Grid
              item
              xl={3}
              md={6}
              xs={12}
            >
              <Mounts mounts={secret.mounted_functions} />
            </Grid>
            <Grid
              item
              xl={6}
              xs={12}
            >
              <Fields secret={secret} disabledEditing={disabledEditing} updateSecret={updateSecret} />
            </Grid>
          </Grid>
        </Box>
      </Container>
    </Page>
  );
}

export default OrderDetailsView;

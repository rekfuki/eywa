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
import CreateSecretForm from './CreateSecretForm';
import Header from './Header';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  }
}));

const SecretCreateView = () => {
  const classes = useStyles();

  console.log("here")
  return (
    <Page
      className={classes.root}
      title="Create Secret"
    >
      <Container maxWidth={false}>
        <Header />
      </Container>
      <Box mt={3}>
        <Container maxWidth="lg">
          <CreateSecretForm />
        </Container>
      </Box>
    </Page>
  );
};

export default SecretCreateView;

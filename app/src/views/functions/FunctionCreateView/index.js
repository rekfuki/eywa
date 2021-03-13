import React, { useState, useEffect, useCallback } from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import {
  Link as RouterLink,
  useHistory
} from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {
  Avatar,
  Box,
  Breadcrumbs,
  Button,
  Card,
  CardContent,
  Container,
  Grid,
  Link,
  Paper,
  Step,
  StepConnector,
  StepLabel,
  Stepper,
  Typography,
  colors,
  makeStyles,
  withStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import {
  Settings as SettingsIcon,
  List as ListIcon,
  Star as StarIcon,
  Sliders as SlidersIcon,
} from 'react-feather';
import Page from 'src/components/Page';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import axios from 'src/utils/axios';
import FunctionDetails from './FunctionDetails';
import EnvVars from './EnvVars';
import Resources from './Resources';

const steps = [
  {
    label: 'Function Details & Secrets',
    icon: SettingsIcon
  },
  {
    label: 'Environment Variables',
    icon: ListIcon
  },
  {
    label: 'Resources',
    icon: SlidersIcon
  }
];

const CustomStepConnector = withStyles((theme) => ({
  vertical: {
    marginLeft: 19,
    padding: 0,
  },
  line: {
    borderColor: theme.palette.divider
  }
}))(StepConnector);

const useCustomStepIconStyles = makeStyles((theme) => ({
  root: {},
  active: {
    backgroundColor: theme.palette.secondary.main,
    boxShadow: theme.shadows[10],
    color: theme.palette.secondary.contrastText
  },
  completed: {
    backgroundColor: theme.palette.secondary.main,
    color: theme.palette.secondary.contrastText
  }
}));

const CustomStepIcon = ({ active, completed, icon }) => {
  const classes = useCustomStepIconStyles();

  const Icon = steps[icon - 1].icon;

  return (
    <Avatar
      className={clsx(classes.root, {
        [classes.active]: active,
        [classes.completed]: completed
      })}
    >
      <Icon size="20" />
    </Avatar>
  );
};

CustomStepIcon.propTypes = {
  active: PropTypes.bool,
  completed: PropTypes.bool,
  icon: PropTypes.number.isRequired
};

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
  avatar: {
    backgroundColor: colors.red[600]
  },
  stepper: {
    backgroundColor: 'transparent'
  }
}));

const FunctionCreateView = () => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();
  const [activeStep, setActiveStep] = useState(0);
  const [completed, setCompleted] = useState(false);
  const [secrets, setSecrets] = useState(null);
  const [images, setImages] = useState(null);

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleComplete = async () => {
    try {
      let payload = getPayload();
      console.log(payload)
      payload = {
        ...payload,
        "secrets": payload.secrets.map((secret) => secret.id),
        "image_id": payload.image.id,
        "read_timeout": `${payload.read_timeout}s`,
        "write_timeout": `${payload.write_timeout}s`,
        "resources": {
          "max_cpu": `${payload.resources.max_cpu}m`,
          "min_cpu": `${payload.resources.min_cpu}m`,
          "max_memory": `${payload.resources.max_memory}Mi`,
          "min_memory": `${payload.resources.min_memory}Mi`,
        }
      };
      delete payload["image"];

      const response = await axios.post("/eywa/api/functions", payload)
      enqueueSnackbar('Function updated', {
        variant: 'success',
      });
      history.push("/app/functions/" + response.data.id)
    } catch (err) {
      console.error(err);
      let msg = "Failed to create function"
      if (err.response.status === 400) {
        msg = err.response.data
      }
      enqueueSnackbar(msg, {
        variant: 'error',
      });
    }
  };

  const getSecrets = useCallback(async () => {
    try {
      const response = await axios.get('/eywa/api/secrets');
      if (isMountedRef.current) {
        setSecrets(response.data.objects);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get secrets', {
        variant: 'error',
      });
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
        variant: 'error',
      });
    }
  }, [isMountedRef]);

  const getPayload = () => {
    const value = window.localStorage.getItem("payload");
    return value !== null
      ? JSON.parse(value)
      : {};
  }

  const setPayload = (value) => {
    window.localStorage.setItem("payload", JSON.stringify(value));
  }

  useEffect(() => {
    getSecrets();
    getImages();
    setPayload({});
  }, []);

  if (!secrets || !images) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Project Create"
    >
      <Container maxWidth="lg">
        <Box mb={3}>
          <Breadcrumbs
            separator={<NavigateNextIcon fontSize="small" />}
            aria-label="breadcrumb"
          >
            <Link
              variant="body1"
              color="inherit"
              to="/app"
              component={RouterLink}
            >
              Dashboard
            </Link>
            <Link
              variant="body1"
              color="inherit"
              to="/app/functions"
              component={RouterLink}
            >
              Functions
            </Link>
            <Typography
              variant="body1"
              color="textPrimary"
            >
              Create
            </Typography>
          </Breadcrumbs>
          <Typography
            variant="h3"
            color="textPrimary"
          >
            Create Wizard &amp; Process
          </Typography>
        </Box>
        {!completed &&
          <Paper>
            <Grid container>
              <Grid
                item
                xs={12}
                md={3}
              >
                <Stepper
                  activeStep={activeStep}
                  className={classes.stepper}
                  connector={<CustomStepConnector />}
                  orientation="vertical"
                >
                  {steps.map((step) => (
                    <Step key={step.label}>
                      <StepLabel StepIconComponent={CustomStepIcon}>
                        {step.label}
                      </StepLabel>
                    </Step>
                  ))}
                </Stepper>
              </Grid>
              <Grid
                item
                xs={12}
                md={9}
              >
                <Box p={3}>
                  {activeStep === 0 && (
                    <FunctionDetails
                      payload={getPayload()}
                      setPayload={setPayload}
                      secrets={secrets}
                      images={images}
                      onNext={handleNext}
                    />
                  )}
                  {activeStep === 1 && (
                    <EnvVars
                      payload={getPayload()}
                      setPayload={setPayload}
                      onBack={handleBack}
                      onNext={handleNext}
                    />
                  )}
                  {activeStep === 2 && (
                    <Resources
                      payload={getPayload()}
                      setPayload={setPayload}
                      onBack={handleBack}
                      onComplete={handleComplete}
                    />
                  )}
                </Box>
              </Grid>
            </Grid>
          </Paper>
        }
      </Container>
    </Page>
  );
};

export default FunctionCreateView;

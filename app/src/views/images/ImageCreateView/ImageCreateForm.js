import React from 'react';
import { useHistory } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { useSnackbar } from 'notistack';
import {
  Box,
  Button,
  Card,
  CardContent,
  IconButton,
  Grid,
  Switch,
  TextField,
  Typography,
  makeStyles
} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import { 
  X as XIcon,
} from 'react-feather';
import axios from 'src/utils/axios';

const useStyles = makeStyles(() => ({
  root: {}
}));

const FunctionEditForm = ({
  className,
  fn,
  secrets,
  ...rest
}) => {
  const classes = useStyles();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();

  var envVarsArr = [];
  for (var key in fn.env_vars){
    envVarsArr.push({"key": key, "value": fn.env_vars[key]})
  } 

  const mountedSecrets = fn.secrets || [];

  console.log(envVarsArr)
  return (
  );
};

FunctionEditForm.propTypes = {
  className: PropTypes.string,
  fn: PropTypes.object.isRequired,
  secrets: PropTypes.array.isRequired
};

export default FunctionEditForm;

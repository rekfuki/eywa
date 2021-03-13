import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import * as Yup from 'yup';
import { Formik } from 'formik';
import {
  Box,
  Button,
  IconButton,
  Grid,
  TextField,
  Typography,
  makeStyles
} from '@material-ui/core';
import {
  X as XIcon,
} from 'react-feather';

const useStyles = makeStyles((theme) => ({
  root: {}
}));

const EnvVars = ({
  className,
  payload,
  setPayload,
  onNext,
  onBack,
  ...rest
}) => {
  const classes = useStyles();

  var envVarsArr = [];
  for (var key in payload.env_vars) {
    envVarsArr.push({ "key": key, "value": payload.env_vars[key] })
  }

  return (
    <Formik
      initialValues={{
        envVars: envVarsArr || [],
        submit: null
      }}
      validationSchema={Yup.object().shape({
        envVars: Yup.array().of(
          Yup.object().shape({
            "key": Yup
              .string()
              .required("Required")
              .min(1, "Must be at least one character long")
              .max(255, "Must be at most 255 characters long"),
            "value": Yup
              .string()
              .required("Required")
              .min(1, "Must be at least one character long")
              .max(2000, "Must be at most 2000 characters long"),
          }))
      })}
      onSubmit={async (values, {
        setErrors,
        setStatus,
        setSubmitting
      }) => {
        try {

          let envVars = {};
          values.envVars.map((obj, index) => {
            envVars[obj.key] = obj.value
          })

          setPayload({
            ...payload,
            env_vars: envVars,
          })
          setStatus({ success: true });
          setSubmitting(false);
          if (onNext) {
            onNext();
          }
        } catch (err) {
          console.error(err);
          setStatus({ success: false });
          setErrors({ submit: err.message });
          setSubmitting(false);
        }
      }}
    >
      {({
        errors,
        handleBlur,
        handleChange,
        handleSubmit,
        setFieldValue,
        isSubmitting,
        touched,
        values
      }) => (
        <form
          onSubmit={handleSubmit}
          className={clsx(classes.root, className)}
          {...rest}
        >
          <Box mb={3} mt={3}>
            <Typography variant="h4" component="h4" gutterBottom>
              Environment Variables
            </Typography>
          </Box>
          <Grid
            container
            direction="column"
            justify="flex-start"
            alignItems="stretch"
          >
            {values.envVars.map((obj, index) => (
              <Box display="flex" flexGrow={1} mb={3} key={index}>
                <Grid
                  container
                  spacing={3}
                >
                  <Grid
                    item
                    md={6}
                    xs={12}
                  >
                    <TextField
                      error={Boolean(
                        (
                          typeof touched.envVars !== 'undefined'
                          &&
                          typeof touched.envVars[index] !== 'undefined'
                          &&
                          touched.envVars[index].key
                        )
                        &&
                        (
                          typeof errors.envVars !== 'undefined'
                          &&
                          typeof errors.envVars[index] !== 'undefined'
                          &&
                          errors.envVars[index].key
                        )
                      )}
                      helperText={
                        (
                          typeof touched.envVars !== 'undefined'
                          &&
                          typeof touched.envVars[index] !== 'undefined'
                          &&
                          touched.envVars[index].key
                        )
                        &&
                        (
                          typeof errors.envVars !== 'undefined'
                          &&
                          typeof errors.envVars[index] !== 'undefined'
                          &&
                          errors.envVars[index].key
                        )
                      }
                      fullWidth
                      label="Key"
                      name={`envVars.${index}.key`}
                      onBlur={handleBlur}
                      onChange={handleChange}
                      value={obj.key}
                      variant="outlined"
                    />
                  </Grid>
                  <Grid
                    item
                    md={5}
                    xs={12}
                  >
                    <TextField
                      error={Boolean(
                        (
                          typeof touched.envVars !== 'undefined'
                          &&
                          typeof touched.envVars[index] !== 'undefined'
                          &&
                          touched.envVars[index].value
                        )
                        &&
                        (
                          typeof errors.envVars !== 'undefined'
                          &&
                          typeof errors.envVars[index] !== 'undefined'
                          &&
                          errors.envVars[index].value
                        )
                      )}
                      helperText={
                        (
                          typeof touched.envVars !== 'undefined'
                          &&
                          typeof touched.envVars[index] !== 'undefined'
                          &&
                          touched.envVars[index].value
                        )
                        &&
                        (
                          typeof errors.envVars !== 'undefined'
                          &&
                          typeof errors.envVars[index] !== 'undefined'
                          &&
                          errors.envVars[index].value
                        )
                      }
                      fullWidth
                      label="Value"
                      name={`envVars.${index}.value`}
                      onBlur={handleBlur}
                      onChange={handleChange}
                      value={obj.value}
                      variant="outlined"
                    />
                  </Grid>
                  <Grid
                    item
                    md={1}
                    xs={12}
                  >
                    <IconButton color="primary" component="span" onClick={() => {
                      if (touched.envVars) { touched.envVars.splice(index, 1) }
                      values.envVars.splice(index, 1)
                      setFieldValue(`envVars`, values.envVars)
                    }}>
                      <XIcon />
                    </IconButton>
                  </Grid>
                </Grid>
              </Box>
            ))}
            <Grid item md={1}>
              <Button variant="contained" color="secondary" onClick={() => {
                values.envVars.push({ "key": "", "value": "" })
                setFieldValue("envVars", values.envVars)
              }}>
                Add
              </Button>
            </Grid>
          </Grid>
          <Box
            mt={6}
            display="flex"
          >
            {onBack && (
              <Button
                onClick={() => {
                  let envVars = {};
                  values.envVars.map((obj, index) => {
                    envVars[obj.key] = obj.value
                  });

                  setPayload({
                    ...payload,
                    env_vars: envVars,
                  });
                  onBack();
                }}
                size="large"
              >
                Previous
              </Button>
            )}
            <Box flexGrow={1} />
            <Button
              color="secondary"
              disabled={isSubmitting}
              type="submit"
              variant="contained"
              size="large"
            >
              Next
            </Button>
          </Box>
        </form>
      )}
    </Formik>
  );
};

EnvVars.propTypes = {
  className: PropTypes.string,
  payload: PropTypes.object.isRequired,
  setPayload: PropTypes.func.isRequired,
  onNext: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default EnvVars;

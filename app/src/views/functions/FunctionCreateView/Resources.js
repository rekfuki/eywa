import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import * as Yup from 'yup';
import { Formik } from 'formik';
import {
  Box,
  Button,
  Grid,
  TextField,
  Typography,
  makeStyles
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {}
}));

const Resources = ({
  className,
  payload,
  setPayload,
  onComplete,
  onBack,
  ...rest
}) => {
  const classes = useStyles();

  return (
    <Formik
      initialValues={{
        minCpu: payload.resources && payload.resources.min_cpu.replace('m', '') || '20',
        maxCpu: payload.resources && payload.resources.max_cpu.replace('m', '') || '50',
        minMemory: payload.resources && payload.resources.min_memory.replace('Mi', '') || '20',
        maxMemory: payload.resources && payload.resources.max_memory.replace('Mi', '') || '50',
      }}
      validationSchema={Yup.object().shape({
        minCpu: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(20, "Must be greater than or equal to 20 milicores")
          .max(500, "Must be less than or equal to 500 milicores")
          .max(Yup.ref("maxCpu"), "Must be less than or equal to Maximum CPU"),
        maxCpu: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(20, "Must be greater than or equal to 20 milicores")
          .min(Yup.ref("minCpu"), "Must be greater than or equal to Minimum CPU")
          .max(500, "Must be less than or equal to 500 milicores"),
        minMemory: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(20, "Must be greater than or equal to 20 mebibytes")
          .max(2000, "Must be less than or equal to 2000 mebibytes")
          .max(Yup.ref("maxMemory"), "Must be less than or equal to Maximum Memory"),
        maxMemory: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(20, "Must be greater than or equal to 20 mebibytes")
          .min(Yup.ref("minMemory"), "Must be greater than or equal to Minimum Memory")
          .max(2000, "Must be less than or equal to 2000 mebibytes"),
      })}
      onSubmit={async (values, {
        setErrors,
        setStatus,
        setSubmitting
      }) => {
        try {

          setPayload({
            ...payload,
            resources: {
              min_cpu: values.minCpu,
              max_cpu: values.maxCpu,
              min_memory: values.minMemory,
              max_memory: values.maxMemory,
            }
          })

          if (onComplete) {
            onComplete();
          }

          setStatus({ success: true });
          setSubmitting(false);
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
              Resources
            </Typography>
          </Box>
          <Grid container spacing={3}>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.minCpu && errors.minCpu)}
                fullWidth
                helperText={touched.minCpu && errors.minCpu}
                label="Request CPU (milicores)"
                name="minCpu"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.minCpu}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.maxCpu && errors.maxCpu)}
                fullWidth
                helperText={touched.maxCpu && errors.maxCpu}
                label="Limit CPU (milicores)"
                name="maxCpu"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.maxCpu}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.minMemory && errors.minMemory)}
                fullWidth
                helperText={touched.minMemory && errors.minMemory}
                label="Request Memory (mebibytes)"
                name="minMemory"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.minMemory}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.maxMemory && errors.maxMemory)}
                fullWidth
                helperText={touched.maxMemory && errors.maxMemory}
                label="Limit Memory (mebibytes)"
                name="maxMemory"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.maxMemory}
                variant="outlined"
              />
            </Grid>
          </Grid>
          <Box
            mt={6}
            display="flex"
          >
            {onBack && (
              <Button
                onClick={() => {
                  setPayload({
                    ...payload,
                    resources: {
                      min_cpu: values.minCpu,
                      max_cpu: values.maxCpu,
                      min_memory: values.minMemory,
                      max_memory: values.maxMemory,
                    }
                  })
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
              Create
            </Button>
          </Box>
        </form>
      )}
    </Formik>
  );
};

Resources.propTypes = {
  className: PropTypes.string,
  payload: PropTypes.object.isRequired,
  setPayload: PropTypes.func.isRequired,
  onComplete: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default Resources;

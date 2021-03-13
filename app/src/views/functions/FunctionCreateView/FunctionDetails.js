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
import Autocomplete from '@material-ui/lab/Autocomplete';

const useStyles = makeStyles((theme) => ({
  root: {}
}));

const FunctionDetails = ({
  className,
  payload,
  setPayload,
  secrets,
  images,
  onNext,
  ...rest
}) => {
  const classes = useStyles();

  const mountedSecrets = payload.secrets
  const setImage = payload.image;

  return (
    <Formik
      initialValues={{
        name: payload.name || '',
        minReplicas: payload.min_replicas || 0,
        maxReplicas: payload.max_replicas || 1,
        scalingFactor: payload.scaling_factor || 20,
        maxConcurrency: payload.max_concurrency || 0,
        readTimeout: payload.read_timeout || '10',
        writeTimeout: payload.write_timeout || '10',
        secrets: payload.secrets || [],
        image: payload.image || {},
        submit: null
      }}
      validationSchema={Yup.object().shape({
        name: Yup
          .string("Must be a string")
          .required("Required")
          .min(5, "Must be at least 5 characters long")
          .matches(/^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$/, "Must match ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"),
        image: Yup.object().shape({
          id: Yup.string().required(),
        }),
        minReplicas: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(0, "Must be greater than or equal to 0")
          .max(Yup.ref("maxReplicas"), "Cannot be great than Maximum Replicas"),
        maxReplicas: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(1, "Must be greater than or equal to 1")
          .min(Yup.ref("minReplicas"), "Must be greater than or equal to Minimum Replicas")
          .max(100, "Must be less than or equal to 100"),
        scalingFactor: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(0, "Must be greater than or equal to 0")
          .max(100, "Must be less than or equal to 100"),
        maxConcurrency: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(0, "Must be greater than or equal to 0"),
        readTimeout: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(0, "Must be greater than or equal to 0"),
        writeTimeout: Yup
          .number()
          .typeError("Must be a number")
          .integer("Must be an integer")
          .required("Required")
          .min(0, "Must be greater than or equal to 0"),
      })}
      onSubmit={async (values, {
        setErrors,
        setStatus,
        setSubmitting
      }) => {
        console.log(values)
        try {
          setPayload({
            ...payload,
            name: values.name,
            image: values.image,
            min_replicas: parseInt(values.minReplicas, 10),
            max_replicas: parseInt(values.maxReplicas, 10),
            scaling_factor: parseInt(values.scalingFactor, 10),
            max_concurrency: parseInt(values.maxConcurrency, 10),
            read_timeout: values.readTimeout,
            write_timeout: values.writeTimeout,
            secrets: values.secrets,
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
          <Box mb={3}>
            <Typography
              variant="h3"
              color="textPrimary"
            >
              Function Details
            </Typography>
          </Box>
          <Grid
            container
            spacing={3}
          >
            <Grid
              item
              md={12}
              xs={12}
            >
              <TextField
                error={Boolean(touched.name && errors.name)}
                fullWidth
                helperText={touched.name && errors.name}
                label="Function Name"
                name="name"
                onBlur={handleBlur}
                onChange={handleChange}
                required
                value={values.name}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={12}
              xs={12}
            >
              <Autocomplete
                name="image"
                options={images}
                getOptionLabel={(option) => `${option.name} (${option.language}) (${option.version})`}
                defaultValue={setImage}
                filterSelectedOptions
                required
                getOptionSelected={(option, value) => option.id === value.id}
                onChange={(_, value) => setFieldValue("image", value)}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    variant="outlined"
                    label="Image"
                    placeholder="Image"
                    required
                  />
                )}
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.minReplicas && errors.minReplicas)}
                fullWidth
                helperText={touched.minReplicas && errors.minReplicas}
                label="Mininum Replicas"
                name="minReplicas"
                onBlur={handleBlur}
                onChange={handleChange}
                required
                value={values.minReplicas}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.maxReplicas && errors.maxReplicas)}
                fullWidth
                helperText={touched.maxReplicas && errors.maxReplicas}
                label="Maximum Replicas"
                name="maxReplicas"
                onBlur={handleBlur}
                onChange={handleChange}
                required
                value={values.maxReplicas}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.scalingFactor && errors.scalingFactor)}
                fullWidth
                helperText={touched.scalingFactor && errors.scalingFactor}
                label="Scaling Factor (%)"
                name="scalingFactor"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.scalingFactor}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.maxConcurrency && errors.maxConcurrency)}
                fullWidth
                helperText={touched.maxConcurrency && errors.maxConcurrency}
                label="Maximum Concurrency"
                name="maxConcurrency"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.maxConcurrency}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.readTimeout && errors.readTimeout)}
                fullWidth
                helperText={touched.readTimeout && errors.readTimeout}
                label="Read Timeout (Seconds)"
                name="readTimeout"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.readTimeout}
                variant="outlined"
              />
            </Grid>
            <Grid
              item
              md={6}
              xs={12}
            >
              <TextField
                error={Boolean(touched.writeTimeout && errors.writeTimeout)}
                fullWidth
                helperText={touched.writeTimeout && errors.writeTimeout}
                label="Write Timeout (Seconds)"
                name="writeTimeout"
                onBlur={handleBlur}
                onChange={handleChange}
                value={values.writeTimeout}
                variant="outlined"
              />
            </Grid>
          </Grid>
          <Box mb={3} mt={3}>
            <Typography variant="h4" component="h4" gutterBottom>
              Mounted Secrets
            </Typography>
          </Box>
          <Grid container spacing={3}>
            <Grid
              item
              md={12}
              xs={12}
            >
              <Autocomplete
                multiple
                name="secrets"
                id="tags-outlined"
                options={secrets}
                getOptionLabel={(option) => option.name}
                defaultValue={mountedSecrets}
                filterSelectedOptions
                getOptionSelected={(option, value) => option.id === value.id}
                onChange={(_, value) => setFieldValue("secrets", value)}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    variant="outlined"
                    label="Attached Secrets"
                    placeholder="Secrets"
                  />
                )}
              />
            </Grid>
          </Grid>
          <Box
            mt={6}
            display="flex"
            justifyContent="flex-end"
          >
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

FunctionDetails.propTypes = {
  className: PropTypes.string,
  payload: PropTypes.object.isRequired,
  setPayload: PropTypes.func.isRequired,
  secrets: PropTypes.array.isRequired,
  images: PropTypes.array.isRequired,
  onNext: PropTypes.func.isRequired,
};

export default FunctionDetails;

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
  X as XIcon
} from 'react-feather';
import axios from 'src/utils/axios';

const useStyles = makeStyles(() => ({
  root: {}
}));

const FunctionEditForm = ({
  className,
  fn,
  secrets,
  images,
  ...rest
}) => {
  const classes = useStyles();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();

  var envVarsArr = [];
  for (var key in fn.env_vars) {
    envVarsArr.push({ "key": key, "value": fn.env_vars[key] })
  }

  const mountedSecrets = fn.secrets || [];
  const setImage = images.filter(image => {
    return image.id === fn.image_id;
  })[0]

  return (
    <Formik
      initialValues={{
        minReplicas: fn.min_replicas || 0,
        maxReplicas: fn.max_replicas || 1,
        scalingFactor: fn.scaling_factor || 20,
        maxConcurrency: fn.max_concurrency || 0,
        readTimeout: fn.read_timeout.replace('s', '') || '10',
        writeTimeout: fn.write_timeout.replace('s', '') || '10',
        // minCpu: fn.resources.min_cpu.replace('m', '') || '20',
        // maxCpu: fn.resources.max_cpu.replace('m', '') || '50',
        // minMemory: fn.resources.min_memory.replace('Mi', '') || '20',
        // maxMemory: fn.resources.max_memory.replace('Mi', '') || '50',
        writeDebug: fn.write_debug || false,
        secrets: fn.secrets || [],
        envVars: envVarsArr || [],
        image: setImage || {},
        submit: null
      }}
      validationSchema={Yup.object().shape({
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
        // minCpu: Yup
        //   .number()
        //   .typeError("Must be a number")
        //   .integer("Must be an integer")
        //   .required("Required")
        //   .min(20, "Must be greater than or equal to 20 milicores")
        //   .max(500, "Must be less than or equal to 500 milicores")
        //   .max(Yup.ref("maxCpu"), "Must be less than or equal to Maximum CPU"),
        // maxCpu: Yup
        //   .number()
        //   .typeError("Must be a number")
        //   .integer("Must be an integer")
        //   .required("Required")
        //   .min(20, "Must be greater than or equal to 20 milicores")
        //   .min(Yup.ref("minCpu"), "Must be greater than or equal to Minimum CPU")
        //   .max(500, "Must be less than or equal to 500 milicores"),
        // minMemory: Yup
        //   .number()
        //   .typeError("Must be a number")
        //   .integer("Must be an integer")
        //   .required("Required")
        //   .min(20, "Must be greater than or equal to 20 mebibytes")
        //   .max(2000, "Must be less than or equal to 2000 mebibytes")
        //   .max(Yup.ref("maxMemory"), "Must be less than or equal to Maximum Memory"),
        // maxMemory: Yup
        //   .number()
        //   .typeError("Must be a number")
        //   .integer("Must be an integer")
        //   .required("Required")
        //   .min(20, "Must be greater than or equal to 20 mebibytes")
        //   .min(Yup.ref("minMemory"), "Must be greater than or equal to Minimum Memory")
        //   .max(2000, "Must be less than or equal to 2000 mebibytes"),
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
              .max(2000, "Must be at most 2000 characters long")
          }))
      })}
      onSubmit={async (values, {
        setErrors,
        setStatus,
        setSubmitting
      }) => {
        try {
          console.log(values)

          const payload = {
            "env_vars": values.envVars.reduce((obj, item) => {
              obj[item.key] = item.value
              return obj
            }, {}),
            "image_id": values.image.id,
            "secrets": values.secrets.map((secret) => secret.id),
            "min_replicas": parseInt(values.minReplicas, 10),
            "max_replicas": parseInt(values.maxReplicas, 10),
            "scaling_factor": parseInt(values.scalingFactor, 10),
            "max_concurrency": parseInt(values.maxConcurrency, 10),
            "write_debug": values.writeDebug,
            "read_timeout": `${values.readTimeout}s`,
            "write_timeout": `${values.writeTimeout}s`
            // "resources": {
            //   "min_cpu": `${values.minCpu}m`,
            //   "max_cpu": `${values.maxCpu}m`,
            //   "min_memory": `${values.minMemory}Mi`,
            //   "max_memory": `${values.maxMemory}Mi`,
            // }
          }

          console.log(payload)

          await axios.put("/eywa/api/functions/" + fn.id, payload);
          setStatus({ success: true });
          setSubmitting(false);
          enqueueSnackbar('Function updated', {
            variant: 'success'
          });
          history.push("/app/functions/" + fn.id)
        } catch (err) {
          console.error(err);
          setStatus({ success: false });
          setErrors({ submit: err.message });
          setSubmitting(false);
          enqueueSnackbar('Failed to update', {
            variant: 'error'
          });
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
          className={clsx(classes.root, className)}
          onSubmit={handleSubmit}
          {...rest}
        >
          <Card>
            <CardContent>
              <Box mb={3}>
                <Typography variant="h4" component="h4" gutterBottom>
                  Function Information
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
                          touched.envVars.splice(index, 1)
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
              {/* <Box mb={3} mt={3}>
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

                <Grid
                  item
                  md={6}
                  xs={12}
                >
                  <Typography component="div">
                    <Typography
                      variant="h5"
                      color="textPrimary"
                    >
                      Write Debug
                    </Typography>
                    <Grid component="label" container alignItems="center" spacing={1}>
                      <Grid item>Off</Grid>
                      <Grid item>
                        <Switch
                          checked={values.writeDebug}
                          color="secondary"
                          edge="start"
                          name="writeDebug"
                          onChange={handleChange}
                          value={values.writeDebug}
                        />
                      </Grid>
                      <Grid item>On</Grid>
                    </Grid>
                  </Typography>
                </Grid>
              </Grid> */}
              <Box mt={2} display="flex" flexDirection="row-reverse">
                <Button
                  variant="contained"
                  color="secondary"
                  type="submit"
                  disabled={isSubmitting}
                >
                  Update
                </Button>
              </Box>
            </CardContent>
          </Card>
        </form>
      )}
    </Formik>
  );
};

FunctionEditForm.propTypes = {
  className: PropTypes.string,
  fn: PropTypes.object.isRequired,
  secrets: PropTypes.array.isRequired
};

export default FunctionEditForm;

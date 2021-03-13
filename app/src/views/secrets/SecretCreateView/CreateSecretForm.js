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
  CardHeader,
  Divider,
  Grid,
  Switch,
  TextField,
  Typography,
  makeStyles
} from '@material-ui/core';
import {
  X as XIcon
} from 'react-feather';
import axios from 'src/utils/axios';

const useStyles = makeStyles(() => ({
  root: {}
}));

const CreateSecretForm = ({
  className,
  ...rest
}) => {
  const classes = useStyles();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();

  return (
    <Formik
      initialValues={{
        name: '',
        fields: [{ key: "", value: "" }],
        submit: null
      }}
      validationSchema={Yup.object().shape({
        name: Yup
          .string("Must be a string")
          .required("Required")
          .min(5, "Must be at least 5 characters long")
          .matches(/^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$/, "Must match ^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"),
        fields: Yup.array().of(
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
          })).min(1, "Must have at least one field set")
      })}
      onSubmit={async (values, {
        setErrors,
        setStatus,
        setSubmitting
      }) => {
        try {
          const payload = {
            "name": values.name,
            "data": values.fields.reduce((obj, item) => {
              obj[item.key] = item.value
              return obj
            }, {})
          }

          const response = await axios.post("/eywa/api/secrets", payload);
          setStatus({ success: true });
          setSubmitting(false);
          enqueueSnackbar('Secret created', {
            variant: 'success'
          });
          history.push("/app/secrets/" + response.data.id)
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
            <CardHeader title="Create Secret" />
            <Divider />
            <CardContent>
              <Grid
                container
                spacing={3}
              >
                <Grid
                  item
                  xs={11}
                >
                  <TextField
                    error={Boolean(touched.name && errors.name)}
                    fullWidth
                    helperText={touched.name && errors.name}
                    label="Secret Name"
                    name="name"
                    onBlur={handleBlur}
                    onChange={handleChange}
                    required
                    value={values.name}
                    variant="outlined"
                  />
                </Grid>
                <Box mb={3} mt={3}>
                </Box>
                <Grid
                  item
                  container
                  direction="column"
                  justify="flex-start"
                  alignItems="stretch"
                >
                  <Typography variant="h5" component="h5" gutterBottom>
                    Secret Fields
                  </Typography>
                  {values.fields.map((obj, index) => (
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
                                typeof touched.fields !== 'undefined'
                                &&
                                typeof touched.fields[index] !== 'undefined'
                                &&
                                touched.fields[index].key
                              )
                              &&
                              (
                                typeof errors.fields !== 'undefined'
                                &&
                                typeof errors.fields[index] !== 'undefined'
                                &&
                                errors.fields[index].key
                              )
                            )}
                            helperText={
                              (
                                typeof touched.fields !== 'undefined'
                                &&
                                typeof touched.fields[index] !== 'undefined'
                                &&
                                touched.fields[index].key
                              )
                              &&
                              (
                                typeof errors.fields !== 'undefined'
                                &&
                                typeof errors.fields[index] !== 'undefined'
                                &&
                                errors.fields[index].key
                              )
                            }
                            fullWidth
                            label="Key"
                            name={`fields.${index}.key`}
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
                                typeof touched.fields !== 'undefined'
                                &&
                                typeof touched.fields[index] !== 'undefined'
                                &&
                                touched.fields[index].value
                              )
                              &&
                              (
                                typeof errors.fields !== 'undefined'
                                &&
                                typeof errors.fields[index] !== 'undefined'
                                &&
                                errors.fields[index].value
                              )
                            )}
                            helperText={
                              (
                                typeof touched.fields !== 'undefined'
                                &&
                                typeof touched.fields[index] !== 'undefined'
                                &&
                                touched.fields[index].value
                              )
                              &&
                              (
                                typeof errors.fields !== 'undefined'
                                &&
                                typeof errors.fields[index] !== 'undefined'
                                &&
                                errors.fields[index].value
                              )
                            }
                            fullWidth
                            label="Value"
                            name={`fields.${index}.value`}
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
                            if (values.fields.length > 1) {
                              touched.fields.splice(index, 1)
                              values.fields.splice(index, 1)
                              setFieldValue(`envVars`, values.fields)
                            }
                          }}>
                            <XIcon />
                          </IconButton>
                        </Grid>
                      </Grid>
                    </Box>
                  ))}
                  <Grid item md={1}>
                    <Button variant="contained" color="secondary" onClick={() => {
                      values.fields.push({ "key": "", "value": "" })
                      setFieldValue("envVars", values.fields)
                    }}>
                      Add
                    </Button>
                  </Grid>
                </Grid>
              </Grid>
              <Box mt={2} display="flex" flexDirection="row-reverse">
                <Button
                  variant="contained"
                  color="secondary"
                  type="submit"
                  disabled={isSubmitting}
                >
                  Create
                </Button>
              </Box>
            </CardContent>
          </Card>
        </form>
      )}
    </Formik >
  );
};

export default CreateSecretForm;

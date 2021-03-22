import React from 'react';
import clsx from 'clsx';
import { Formik } from 'formik';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  IconButton,
  TextField,
  Tooltip,
  Typography,
  useTheme,
  makeStyles
} from '@material-ui/core';
import {
  X as XIcon
} from 'react-feather';

const useStyles = makeStyles((theme) => ({
  root: {},
  fontWeightMedium: {
    fontWeight: theme.typography.fontWeightMedium
  },
  emptyError: {
    color: theme.palette.error.main
  }
}));

const Fields = ({
  secret,
  disabledEditing,
  updateSecret,
  className,
  ...rest
}) => {
  const theme = useTheme();
  const classes = useStyles(theme);

  const fieldsArr = secret.data_fields.map(field => {
    return { key: field, value: undefined }
  })

  console.log("rendering in fields:", secret);
  console.log("fieldsarr:", fieldsArr);

  return (
    <Formik
      enableReinitialize
      initialValues={{
        deletes: [],
        updates: fieldsArr,

        submit: null
      }}
      validate={(values) => {
        let unchangedValues = 0;
        let valuesErrors = {};
        for (let i = 0; i < values.updates.length; i++) {
          const vals = values.updates[i];

          if (vals.value === undefined) {
            unchangedValues++;
            continue;
          }

          let err = {}
          if (vals.key === "") {
            err.key = "Key is required"
          } else if (vals.key.length > 255) {
            err.key = "Key must be at most 255 characters long"
          }

          if (vals.value === "") {
            err.value = "Value cannot bet empty"
          } else if (vals.value.length > 2000) {
            err.key = "Value must be at most 2000 characters long"
          }

          if (Object.keys(err).length > 0) {
            valuesErrors[i] = err;
          }
        }

        let errors = {};
        if (Object.keys(valuesErrors).length > 0) {
          errors["updates"] = valuesErrors;
        }

        if (values.updates.length == 0) {
          errors["empty"] = true;
        }

        if (values.updates.length == unchangedValues && values.deletes.length == 0) {
          errors["unchanged"] = true;
        }

        return errors;
      }}

      onSubmit={async (values) => {
        const payload = {
          upserts: values.updates.reduce((obj, item) => {
            if (item.value !== undefined) {
              obj[item.key] = item.value
            }
            return obj
          }, {}),
          deletes: values.deletes
        }
        updateSecret(payload);
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
            <CardHeader title="Secret Fields" />
            <Divider />
            <CardContent>
              {errors.empty && !errors.unchanged &&
                < Typography className={classes.emptyError} variant="h4" component="h4" gutterBottom>
                  Secret must have at least one field set
                  </Typography>
              }
              <Grid
                container
                direction="column"
                justify="flex-start"
                alignItems="stretch"
              >
                {values.updates.map((obj, index) => (
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
                          fullWidth
                          error={touched?.updates?.[index]?.key && errors?.updates?.[index]?.key !== undefined}
                          helperText={touched?.updates?.[index]?.key && errors?.updates?.[index]?.key}
                          label="Key"
                          name={`updates.${index}.key`}
                          onBlur={handleBlur}
                          onChange={handleChange}
                          disabled={obj.value === undefined}
                          value={obj.key}
                          variant="outlined"
                          type="text"
                          inputProps={{
                            autoComplete: "off"
                          }}
                        />
                      </Grid>
                      <Grid
                        item
                        md={5}
                        xs={12}
                      >
                        <TextField
                          error={touched?.updates?.[index]?.value && errors?.updates?.[index]?.value !== undefined}
                          helperText={touched?.updates?.[index]?.value && errors?.updates?.[index]?.value}
                          fullWidth
                          label="Value"
                          name={`updates.${index}.value`}
                          onBlur={handleBlur}
                          onChange={handleChange}
                          disabled={obj.value === undefined}
                          value={obj.value ?? "***********"}
                          variant="outlined"
                          type="text"
                          inputProps={{
                            autoComplete: "off"
                          }}
                        />
                      </Grid>
                      <Grid
                        item
                        md={1}
                        xs={12}
                      >
                        <IconButton
                          color="primary"
                          component="span"
                          disabled={disabledEditing}
                          onClick={() => {
                            touched?.updates?.splice(index, 1);
                            values.updates.splice(index, 1);
                            setFieldValue("updates", values.updates);
                            if (obj.key !== "") {
                              values.deletes.push(obj.key);
                            }
                            setFieldValue(`deletes`, values.deletes)
                          }}>
                          <XIcon />
                        </IconButton>
                      </Grid>
                    </Grid>
                  </Box>
                ))}
                <Grid item md={1}>
                  <Button
                    variant="contained"
                    color="secondary"
                    disabled={disabledEditing}
                    onClick={() => {
                      values.updates.push({ "key": "", "value": "" })
                      setFieldValue("updates", values.updates)
                    }}>
                    Add
                  </Button>
                </Grid>
              </Grid>
              <Box mt={2} display="flex" flexDirection="row-reverse">
                <Button
                  variant="contained"
                  color="secondary"
                  type="submit"
                  disabled={isSubmitting || (errors.unchanged == true || errors.empty == true) !== (touched.updates === undefined)}
                >
                  Update
                </Button>
              </Box>
            </CardContent>
          </Card>
        </form>
      )
      }
    </Formik >
  );
};

export default Fields;

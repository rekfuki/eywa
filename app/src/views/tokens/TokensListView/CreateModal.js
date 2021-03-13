import React, {
  useState,
  useCallback,
  useEffect
} from 'react';
import { useHistory } from 'react-router-dom';
import * as Yup from 'yup';
import { Formik } from 'formik';
import moment from 'moment';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Container,
  Dialog,
  DialogContent,
  Divider,
  IconButton,
  Grid,
  InputAdornment,
  Tooltip,
  Typography,
  TextField,
  makeStyles
} from '@material-ui/core';
import {
  X as XIcon,
  Clipboard as ClipboardIcon
} from 'react-feather';
import { useSnackbar } from 'notistack';
import axios from 'src/utils/axios';

const useStyles = makeStyles((theme) => ({
  root: {
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
  danger: {
    color: theme.palette.error.main
  },
  actions: {
    marginTop: theme.spacing(2),
    display: 'flex',
    justifyContent: 'flex-end',
    '& > * + *': {
      marginLeft: theme.spacing(2)
    }
  }
}));

const CreateModal = ({
  onClose,
  open
}) => {
  const classes = useStyles();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();
  const [newToken, setNewToken] = useState(null);

  return (
    <Dialog open={open} onClose={() => {
      if (!newToken) {
        onClose(false);
        return
      }

      setNewToken(null);
      onClose(true);
    }} aria-labelledby="form-dialog-title">
      <DialogContent style={{ padding: 0 }}>
        <Box>
          <Container maxWidth="md" style={{ padding: 0 }}>
            <Card>
              <CardHeader title="Token Information" />
              <Divider />
              <CardContent>
                <Formik
                  initialValues={{
                    name: "",
                    expiresAt: "",
                    submit: null
                  }}
                  validationSchema={Yup.object().shape({
                    name: Yup
                      .string()
                      .typeError("Must be a string")
                      .required("Required")
                      .min(5, "Must be greater than or equal to 5 characters")
                      .max(63, "Must be less than or equal to 63 characters"),
                    expiresAt: Yup
                      .string()
                      .nullable()
                      .test("is-greater", "Should be in the future", function (value) {
                        if (value === "") {
                          return true;
                        }
                        return moment(value).isSameOrAfter(moment());
                      })
                  })}
                  onSubmit={async (values, {
                    setErrors,
                    setStatus,
                    setSubmitting
                  }) => {
                    try {

                      const payload = {
                        name: values.name
                      };

                      if (values.expiresAt !== "") {
                        payload.expires_at = moment(values.expiresAt).valueOf()
                      }

                      const response = await axios.post("/eywa/api/tokens", payload)
                      console.log(response)
                      setNewToken(response.data)

                      setStatus({ success: true });
                      setSubmitting(false);
                      enqueueSnackbar('Token created', {
                        variant: 'success'
                      });

                    } catch (err) {
                      console.error(err);
                      setStatus({ success: false });
                      setErrors({ submit: err.message });
                      setSubmitting(false);

                      enqueueSnackbar("Failed to create token", {
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
                      className={classes.root}
                      onSubmit={handleSubmit}
                    >
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
                            label="Name of the Token"
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
                          md={6}
                          xs={12}
                        >
                          <TextField
                            error={Boolean(touched.expiresAt && errors.expiresAt)}
                            helperText={touched.expiresAt && errors.expiresAt}
                            id="datetime-local"
                            label="Token expiry date"
                            type="datetime-local"
                            name="expiresAt"
                            onBlur={handleBlur}
                            value={values.expiresAt}
                            onChange={handleChange}
                            className={classes.textField}
                            InputProps={{
                              endAdornment: values.expiresAt == "" ? undefined : (
                                <InputAdornment position="end">
                                  <IconButton onClick={() => setFieldValue("expiresAt", "")} size="small" style={{ borderRadius: 0 }} >
                                    <XIcon style={{ paddingBottom: "2px" }} />
                                  </IconButton>
                                </InputAdornment>
                              )
                            }}
                            InputLabelProps={{
                              shrink: true
                            }}
                          />
                        </Grid>
                      </Grid>
                      {!newToken && <Box mt={2} display="flex" flexDirection="row-reverse">
                        <Button
                          variant="contained"
                          color="secondary"
                          type="submit"
                          disabled={isSubmitting}
                        >
                          Create
                      </Button>
                      </Box>}
                    </form>
                  )}
                </Formik>
                {newToken &&
                  <Box mt={3} className={classes.root}>
                    <Typography className={classes.danger}>
                      Please copy this token now as its value will become unavailable after you leave this page
                  </Typography>
                    <TextField
                      id="standard-full-width"
                      label="New token"
                      style={{ margin: 8 }}
                      placeholder="Placeholder"
                      fullWidth
                      value={newToken.token}
                      margin="normal"
                      type="password"
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Tooltip title="Copy to clipboard">
                              <IconButton onClick={() => {
                                navigator.clipboard.writeText(newToken.token)
                                enqueueSnackbar("Copied!", {
                                  variant: 'info'
                                });
                              }} size="small" style={{ borderRadius: 0 }} >
                                <ClipboardIcon style={{ paddingBottom: "2px" }} />
                              </IconButton>
                            </Tooltip>
                          </InputAdornment>
                        )
                      }}
                      InputLabelProps={{
                        shrink: true
                      }}
                    />
                  </Box>
                }
              </CardContent>
            </Card>
          </Container>
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default CreateModal;

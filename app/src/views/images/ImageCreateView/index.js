import React, {
  useState,
  useCallback
} from 'react';
import { useHistory } from 'react-router-dom';
import { useDropzone } from 'react-dropzone';
import * as Yup from 'yup';
import { Formik } from 'formik';
import clsx from 'clsx';
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  Grid,
  Link,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Typography,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  makeStyles
} from '@material-ui/core';
import { useSnackbar } from 'notistack';
import axios from 'src/utils/axios';
import FileCopyIcon from '@material-ui/icons/FileCopy';
import Page from 'src/components/Page';
import Header from './Header';
import bytesToSize from 'src/utils/bytesToSize';

const FilesDropzone = ({ setFieldValue, ...rest }) => {
  const classes = useStyles();
  const [files, setFiles] = useState([]);

  const handleDrop = useCallback((acceptedFiles) => {
    setFiles(acceptedFiles);
    setFieldValue("files", acceptedFiles)
  }, []);

  const handleRemoveAll = () => {
    setFiles([]);
    setFieldValue("files", [])
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    maxFiles: 1,
    accept: '.zip',
    onDrop: handleDrop
  });

  return (
    <div
      className={classes.root}
      {...rest}
    >
      <div
        className={clsx({
          [classes.dropZone]: true,
          [classes.dragActive]: isDragActive
        })}
        {...getRootProps()}
      >
        <input {...getInputProps()} />
        <div>
          <img
            alt="Select file"
            className={classes.image}
            src="/static/images/undraw_add_file2_gvbb.svg"
          />
        </div>
        <div>
          <Typography
            gutterBottom
            variant="h3"
          >
            Select a ZIP fIle
          </Typography>
          <Box mt={2}>
            <Typography
              color="textPrimary"
              variant="body1"
            >
              Drop ZIP archive here or click
              {' '}
              <Link underline="always">browse</Link>
              {' '}
              thorough your machine
            </Typography>
          </Box>
        </div>
      </div>
      {files.length > 0 && (
        <>
          <List className={classes.list}>
            {files.map((file, i) => (
              <ListItem
                divider={i < files.length - 1}
                key={i}
              >
                <ListItemIcon>
                  <FileCopyIcon />
                </ListItemIcon>
                <ListItemText
                  primary={file.name}
                  primaryTypographyProps={{ variant: 'h5' }}
                  secondary={bytesToSize(file.size)}
                />
              </ListItem>
            ))}
          </List>
          <div className={classes.actions}>
            <Button
              onClick={handleRemoveAll}
              size="small"
            >
              Remove
            </Button>
          </div>
        </>
      )}
    </div>
  );
};

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
  dropZone: {
    border: `1px dashed ${theme.palette.divider}`,
    padding: theme.spacing(6),
    outline: 'none',
    display: 'flex',
    justifyContent: 'center',
    flexWrap: 'wrap',
    alignItems: 'center',
    '&:hover': {
      backgroundColor: theme.palette.action.hover,
      opacity: 0.5,
      cursor: 'pointer'
    }
  },
  dragActive: {
    backgroundColor: theme.palette.action.active,
    opacity: 0.5
  },
  image: {
    width: 130
  },
  info: {
    marginTop: theme.spacing(1)
  },
  list: {
    maxHeight: 320
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

const ImageCreateView = () => {
  const classes = useStyles();
  const history = useHistory();
  const { enqueueSnackbar } = useSnackbar();

  return (
    <Page
      className={classes.root}
      title="Create Image"
    >
      <Container maxWidth={false}>
        <Header />
      </Container>
      <Box mt={3}>
        <Container maxWidth="md">
          <Formik
            initialValues={{
              name: "",
              language: "go",
              version: "1.0.0",
              executablePath: "",
              files: [],
              submit: null
            }}
            validationSchema={Yup.object().shape({
              name: Yup
                .string()
                .typeError("Must be a string")
                .required("Required")
                .min(5, "Must be greater than or equal to 5 characters")
                .max(63, "Must be less than or equal to 63 characters"),
              version: Yup
                .string()
                .typeError("Must be a string")
                .required("Required")
                .matches(/^(\d{1,3}\.?){3}$/, "Must match normal SemVer 2.0"),
              language: Yup
                .string()
                .typeError("Must be a string")
                .required("Required"),
              executablePath: Yup
                .string()
                .typeError("Must be a string")
                .min(1, "Must be at least 1 characters long")
                .matches(/^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\/[a-z0-9]([a-z0-9-]*[a-z0-9])?)*(\/[a-z0-9]([a-z0-9-\.]*[a-z0-9])?)?$/,
                  "Path must be a valid relative path (i.e. foo | foo/bar | foo/bar.sh )"),
              files: Yup.array().min(1, "File is required")
            })}
            onSubmit={async (values, {
              setErrors,
              setStatus,
              setSubmitting
            }) => {
              try {
                let formData = new FormData();
                formData.append("source", values.files[0])
                formData.append("version", values.version)
                formData.append("name", values.name)
                formData.append("language", values.language)

                if (values.language === "custom") {
                  formData.append("executable_path", values.executablePath)
                }

                const response = await axios.post("/eywa/api/images", formData, {
                  headers: {
                    'Content-Type': `multipart/form-data; boundary=${formData._boundary}`
                  }
                })

                const buildID = response.data.build_id

                setStatus({ success: true });
                setSubmitting(false);
                enqueueSnackbar('Image created', {
                  variant: 'success'
                });

                history.push(`/app/images/${buildID}/buildlogs`)

              } catch (err) {
                console.error(err);
                setStatus({ success: false });
                setErrors({ submit: err.message });
                setSubmitting(false);

                let message = 'Failed to create image';
                if (err.response.status === 409) {
                  message = "Image with the same name already exists"
                }
                enqueueSnackbar(message, {
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
                <Card>
                  <CardContent>
                    <Box mb={3}>
                      <Typography variant="h4" component="h4" gutterBottom>
                        Image Information
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
                          label="Name of the image"
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
                          error={Boolean(touched.version && errors.version)}
                          fullWidth
                          helperText={touched.version && errors.version}
                          label="Version (SemVer 2.0 normal)"
                          name="version"
                          onBlur={handleBlur}
                          onChange={handleChange}
                          required
                          value={values.version}
                          variant="outlined"
                        />
                      </Grid>
                      <Grid
                        item
                        md={6}
                        xs={12}
                      >
                        <FormControl style={{ display: "flex" }} className={classes.formControl}>
                          <InputLabel id="language-select">Language</InputLabel>
                          <Select
                            defaultValue={values.language}
                            labelId="language-select"
                            name="language"
                            value={values.language}
                            onChange={handleChange}
                            onBlur={handleBlur}
                            variant="outlined"
                          >
                            <MenuItem value={"go"}>Go</MenuItem>
                            <MenuItem value={"node14"}>Node 14</MenuItem>
                            <MenuItem value={"python3"}>Python 3</MenuItem>
                            <MenuItem value={"ruby"}>Ruby</MenuItem>
                            <MenuItem value={"csharp"}>C#</MenuItem>
                            <MenuItem value={"custom"}>Custom</MenuItem>
                          </Select>
                        </FormControl>
                      </Grid>
                      {values.language === "custom" &&
                        <Grid
                          item
                          xs={12}
                        >
                          <TextField
                            error={Boolean(touched.executablePath && errors.executablePath)}
                            fullWidth
                            helperText={touched.executablePath && errors.executablePath}
                            label="Executable filepath (relative to your zip structure)"
                            name="executablePath"
                            onBlur={handleBlur}
                            onChange={handleChange}
                            required
                            value={values.executablePath}
                            variant="outlined"
                          />
                        </Grid>
                      }
                    </Grid>
                    <Grid item xs={12} md={12}>
                      <FilesDropzone
                        setFieldValue={setFieldValue}
                      />
                    </Grid>
                    <Box mt={2} display="flex" flexDirection="row-reverse">
                      <Button
                        variant="contained"
                        color="secondary"
                        type="submit"
                        disabled={isSubmitting || values.files.length == 0}
                      >
                        Create
                      </Button>
                    </Box>
                  </CardContent>
                </Card>
              </form>
            )}
          </Formik>
        </Container>
      </Box>
    </Page>
  );
};

export default ImageCreateView;

import React, {
  useState,
  useEffect,
  useRef,
  useCallback,
} from 'react';
import { useParams } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Box,
  CardContent,
  Container,
  Divider,
  Typography,
  Paper,
  makeStyles
} from '@material-ui/core';
import LinearProgress from '@material-ui/core/LinearProgress';
import _ from "lodash";
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import Label from 'src/components/Label';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Header from './Header';


const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
  valueContainer: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
}));

const ImageBuildView = () => {
  const { imageId } = useParams();
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const fieldRef = useRef(null);
  const { enqueueSnackbar } = useSnackbar();
  const [buildInfo, setBuildInfo] = useState([]);
  const [image, setImage] = useState(null)
  const [building, setBuilding] = useState(false);
  const [userScrolled, setUserScrolled] = useState(false);

  const getImage = async () => {
    try {
      const url = `/eywa/api/images/${imageId}`
      const response = await axios.get(url)

      const data = response.data

      setImage(data)
      setBuilding(data.state === 'building')
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get image', {
        variant: 'error',
      });
    }
  };

  const getImageBuildInfo = async () => {
    try {
      const url = `/eywa/api/images/${imageId}/buildlogs`
      const response = await fetch(url);
      const reader = response.body
        .pipeThrough(new TextDecoderStream())
        .getReader();

      while (true) {
        const { value, done } = await reader.read();
        if (done) {
          getImage()
          break;
        }

        setBuildInfo(buildInfo => [...buildInfo, value]);
      }
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get images build logs', {
        variant: 'error',
      });
    }
  };

  const handleScroll = () => {
    setUserScrolled(true)
  };

  useEffect(() => {
    window.addEventListener('wheel', handleScroll, { passive: true });
    return () => window.removeEventListener('wheel', handleScroll);
  }, [])

  useEffect(() => {
    getImage();
    getImageBuildInfo()

  }, []);

  useEffect(() => {
    if (fieldRef.current && !userScrolled) {
      fieldRef.current.scrollIntoView({ block: "end" });
    }
  }, [buildInfo, image])

  if (!image) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Timeline Details"
    >
      <Container maxWidth={false} onScroll={handleScroll}>
        <Header imageId={imageId} />
        <Box mt={3}>
          <Paper>
            <CardContent>
              <Box display="flex">
                <Box mr={1}>
                  <Typography variant="body1">
                    IMAGE ID:
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="h5">
                    {image.id}
                  </Typography>
                </Box>
                <Box ml={3} mr={1}>
                  <Typography variant="body1">
                    IMAGE NAME:
                    </Typography>
                </Box>
                <Box>
                  <Typography variant="h5">
                    {image.name}
                  </Typography>
                </Box>
                <Box ml={3} mr={1}>
                  <Typography variant="body1">
                    Status:
                    </Typography>
                </Box>
                <Box>
                  <div className={classes.valueContainer}>
                    <Label
                      color={
                        image.state === "building"
                          ? 'warning'
                          : image.state === "success"
                            ? 'success'
                            : 'error'
                      }
                    >
                      {image.state}
                    </Label>
                  </div>
                </Box>
              </Box>
            </CardContent>
            <Divider />
            {building && <LinearProgress />}
            <CardContent>
              <PerfectScrollbar>
                <Typography ref={fieldRef} style={{ whiteSpace: 'pre-line' }} variant="body1">
                  {buildInfo.map((log) => log)}
                </Typography>
              </PerfectScrollbar>
            </CardContent>
          </Paper>
        </Box>
      </Container>
    </Page>
  );
};

export default ImageBuildView;

import React, {
  Fragment,
  useState,
  useEffect,
  useCallback
} from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { useSnackbar } from 'notistack';
import _ from "lodash";
import {
  Box,
  Button,
  Card,
  CardHeader,
  Divider,
  Link,
  LinearProgress,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  makeStyles
} from '@material-ui/core';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import axios from 'src/utils/axios';
import Label from 'src/components/Label';

const useStyles = makeStyles(() => ({
  root: {},
  queryField: {
    width: 500
  }
}));

const Timelines = ({
  className,
  ...rest
}) => {
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const [timelines, setTimelines] = useState({});
  const [loading, setLoading] = useState(true);

  const getTimelines = async () => {
    try {
      let url = `/eywa/api/timeline?page=1&per_page=5`
      const response = await axios.get(url)

      setTimelines(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get timelines', {
        variant: 'error'
      });
    }
    setLoading(false);
  };

  useEffect(() => {
    getTimelines();
  }, []);

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Timeline logs" />
      <Divider />
      <PerfectScrollbar>
        <Box minWidth={700}>
          {loading
            ? <LinearProgress />
            :
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    REQUEST ID
                    </TableCell>
                  <TableCell>
                    FUNCTION ID
                    </TableCell>
                  <TableCell>
                    FUNCTION NAME
                    </TableCell>
                  <TableCell>
                    STATUS
                    </TableCell>
                  <TableCell>
                    DURATION
                    </TableCell>
                  <TableCell>
                    CREATED AT
                    </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {timelines.objects.map((timeline) => (
                  <TableRow hover key={timeline.request_id} onClick={(event) => console.log(event)}>
                    <TableCell>
                      <Link
                        component={RouterLink}
                        to={`/app/timelines/${timeline.request_id}`}
                      >
                        <Typography
                          variant="h6"
                        >
                          {timeline.request_id}
                        </Typography>
                      </Link>
                    </TableCell>
                    <TableCell>
                      <Link
                        component={RouterLink}
                        to={`/app/functions/${timeline.function_id}`}
                      >
                        <Typography
                          variant="h6"
                        >
                          {timeline.function_id}
                        </Typography>
                      </Link>
                    </TableCell>
                    <TableCell>
                      <Typography
                        variant="h6"
                      >
                        {timeline.function_name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Label
                        color={
                          timeline.is_error
                            ? 'error'
                            : 'success'
                        }
                      >
                        {timeline.status}
                      </Label>
                    </TableCell>
                    <TableCell>
                      {
                        moment.duration(timeline.duration).minutes() > 0
                          ?
                          `${moment.duration(timeline.duration).minutes()} min`
                          :
                          moment.duration(timeline.duration).seconds() > 0
                            ?
                            `${moment.duration(timeline.duration).seconds()} s`
                            :
                            `${moment.duration(timeline.duration).milliseconds()} ms`
                      }
                    </TableCell>
                    <TableCell>
                      {moment(timeline.age).format('YYYY/MM/DD | HH:mm:ss')}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          }
        </Box>
      </PerfectScrollbar>
      <Box
        p={2}
        display="flex"
        justifyContent="flex-end"
      >
        <Button
          component={RouterLink}
          size="small"
          to="/app/timelines"
          endIcon={<NavigateNextIcon />}
        >
          See all
        </Button>
      </Box>
    </Card >
  );
};

export default Timelines;

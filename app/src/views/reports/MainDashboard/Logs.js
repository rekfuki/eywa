import React, {
  useState,
  useEffect,
  Fragment
} from 'react';
import { Link as RouterLink } from 'react-router-dom';
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
  Collapse,
  Divider,
  IconButton,
  Typography,
  LinearProgress,
  Link,
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

const useRowStyles = makeStyles({
  root: {
    '& > *': {
      borderBottom: 'unset'
    }
  }
});

function Row(props) {
  const { row } = props;
  const [open, setOpen] = useState(false);
  const classes = useRowStyles();

  return (
    <Fragment>
      <TableRow className={classes.root}>
        <TableCell>
          <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell>
          <Label color={!row.is_error ? 'success' : 'error'}>
            {!row.is_error
              ? 'OK'
              : 'ERROR'
            }
          </Label>
        </TableCell>
        <TableCell>
          <Link
            component={RouterLink}
            to={`/app/functions/${row.function_id}`}
          >
            <Typography
              variant="h6"
            >
              {row.function_id}
            </Typography>
          </Link>
        </TableCell>
        <TableCell>{row.function_name}</TableCell>
        <TableCell>
          <Link
            component={RouterLink}
            to={`/app/timelines/${row.request_id}`}
          >
            <Typography
              variant="h6"
            >
              {row.request_id}
            </Typography>
          </Link>
        </TableCell>
        <TableCell>{row.message}</TableCell>
        <TableCell>
          {moment(row.created_at).format('YYYY/MM/DD | HH:mm:ss')}
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={7}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box margin={1}>
              <Typography variant="h6" gutterBottom component="div">
                Raw Data
              </Typography>
              <Typography component="span" variant="body1">
                <Box bgcolor="secondary" color="primary">
                  <pre style={{ backgroundColor: "inherit" }} bgcolor="primary">
                    {JSON.stringify(row, null, 2)}
                  </pre>
                </Box>
              </Typography>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </Fragment>
  );
}

const useStyles = makeStyles((theme) => ({
  root: {},
  methodCell: {
    width: 100
  },
  statusCell: {
    width: 64
  },
  formControl: {
    minWidth: 120
  }
}));


const Logs = ({ className, ...rest }) => {
  const classes = useStyles();
  const [logs, setLogs] = useState({})
  const [loading, setLoading] = useState(true)
  const { enqueueSnackbar } = useSnackbar();

  const getLogs = async () => {
    try {
      const response = await axios.post(`/eywa/api/events/query?page=1&per_page=5`)

      setLogs(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get logs', {
        variant: 'error'
      });
    }
    setLoading(false);
  };


  useEffect(() => {
    getLogs();
  }, []);

  if (logs.length == 0 && !loading) {
    return null
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Latest Logs" />
      <Divider />
      <PerfectScrollbar>
        <Box minWidth={700}>
          {loading
            ? <LinearProgress />
            : <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                  </TableCell>
                  <TableCell>
                    STATUS
                  </TableCell>
                  <TableCell>
                    FUNCTION ID
                  </TableCell>
                  <TableCell>
                    FUNCTION NAME
                  </TableCell>
                  <TableCell>
                    REQUEST ID
                  </TableCell>
                  <TableCell>
                    MESSAGE
                  </TableCell>
                  <TableCell>
                    TIMESTAMP
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {logs.objects.map((log) => (
                  <Row key={log.id} row={log}></Row>
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
          to="/app/logs"
          endIcon={<NavigateNextIcon />}
        >
          See all
        </Button>
      </Box>
    </Card>
  );
}

export default Logs;

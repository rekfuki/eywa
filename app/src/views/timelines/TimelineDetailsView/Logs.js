import React, {
  useState,
  useEffect,
  useCallback,
  Fragment
} from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { useSnackbar } from 'notistack';
import _ from "lodash";
import {
  Box,
  Card,
  CardHeader,
  Checkbox,
  Collapse,
  Divider,
  FormControlLabel,
  Grid,
  InputAdornment,
  IconButton,
  MenuItem,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TablePagination,
  TextField,
  Select,
  SvgIcon,
  makeStyles
} from '@material-ui/core';
import LinearProgress from '@material-ui/core/LinearProgress';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import {
  Search as SearchIcon
} from 'react-feather';
import axios from 'src/utils/axios';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Label from 'src/components/Label';

const useRowStyles = makeStyles({
  root: {
    '& > *': {
      borderBottom: 'unset',
    },
  },
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
        <TableCell>{row.function_name}</TableCell>
        <TableCell>{row.request_id}</TableCell>
        <TableCell>{row.message}</TableCell>
        <TableCell>
          {moment(row.created_at).format('YYYY/MM/DD | hh:mm:ss')}
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
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
  },
}));


const Logs = ({
  className,
  requestId,
  ...rest }) => {
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const [logs, setLogs] = useState(null);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [showOnlyErrors, setShowOnlyErrors] = useState(false);
  const [startDate, setStartDate] = useState(moment().format("yyyy-MM-DDTHH:mm:ss"));
  const [endDate, setEndDate] = useState(moment().add(-7, 'd').format("yyyy-MM-DDTHH:mm:ss"));
  const [loading, setLoading] = useState(true);

  const delayedQuery = useCallback(_.debounce((q) => {
    getLogs(q, requestId, startDate, endDate, showOnlyErrors, page, limit)
  }, 300), []);

  const getLogs = async (q, requestId, startDate, endDate, showOnlyErrors, page, limit) => {
    setLoading(true);
    try {
      const payload = {
        "request_id": requestId,
        "query": q,
        "timestamp_max": moment(startDate).format("yyyy-MM-DDTHH:mm:ssZ"),
        "timestamp_min": moment(endDate).format("yyyy-MM-DDTHH:mm:ssZ"),
        "only_errors": showOnlyErrors
      }
      const response = await axios.post(`/eywa/api/events/query?page=${page + 1}&per_page=${limit}`, payload)

      setLogs(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get logs', {
        variant: 'error'
      });
    }
    setLoading(false);
  };

  const handleStartDateChange = (event) => {
    setStartDate(event.target.value);
  };

  const handleEndDateChange = (event) => {
    setEndDate(event.target.value);
  };

  const handleShowOnlyErrorsChange = (event) => {
    setShowOnlyErrors(!showOnlyErrors);
  };

  const handleQueryChange = (event) => {
    event.persist();
    setQuery(event.target.value);
    delayedQuery(event.target.value);
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage);
  };

  const handleLimitChange = (event) => {
    setLimit(parseInt(event.target.value));
  };


  useEffect(() => {
    getLogs(query, requestId, startDate, endDate, showOnlyErrors, page, limit);
  }, [page, limit, startDate, endDate, showOnlyErrors]);

  if (!logs && !loading) {
    return null;
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Execution Logs" />
      <Box
        p={2}
        minHeight={56}
        display="flex"
        alignItems="center"
      >
        <Grid
          container
          spacing={3}
        >
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SvgIcon
                      fontSize="small"
                      color="action"
                    >
                      <SearchIcon />
                    </SvgIcon>
                  </InputAdornment>
                )
              }}
              onChange={handleQueryChange}
              placeholder="Search execution logs"
              value={query}
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12} sm={12} md={"auto"}>
            <TextField
              fullWidth
              label="Start Time"
              type="datetime-local"
              defaultValue={startDate}
              className={classes.textField}
              onChange={handleStartDateChange}
              InputLabelProps={{
                shrink: true
              }}
            />
          </Grid>
          <Grid item xs={12} sm={12} md={"auto"}>
            <TextField
              fullWidth
              label="End Time"
              type="datetime-local"
              defaultValue={endDate}
              className={classes.textField}
              onChange={handleEndDateChange}
              InputLabelProps={{
                shrink: true
              }}
            />
          </Grid>
          <Grid item xs={12} sm={12} md={"auto"}>
            <FormControlLabel
              control={
                <Checkbox
                  color="primary"
                  onChange={handleShowOnlyErrorsChange}
                />}
              label="Show only errors"
              labelPlacement="end"
            />
          </Grid>
        </Grid>
      </Box>
      <Divider />
      {loading ? <LinearProgress />
        :
        <div>
          <PerfectScrollbar>
            <Box minWidth={1150}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>
                    </TableCell>
                    <TableCell>
                      STATUS
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
                    <Row key={log.message} row={log}></Row>
                  ))}
                </TableBody>
              </Table>
            </Box>
          </PerfectScrollbar>
          <TablePagination
            component="div"
            count={logs.total_count}
            onChangePage={handlePageChange}
            onChangeRowsPerPage={handleLimitChange}
            page={page}
            rowsPerPage={limit}
            rowsPerPageOptions={[5, 10, 25]}
          />
        </div>}
    </Card>
  );
};

Logs.propTypes = {
  className: PropTypes.string,
  requestId: PropTypes.string.isRequired
};

export default Logs;

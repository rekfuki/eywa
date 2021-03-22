import React, {
  Fragment,
  useState,
  useEffect,
  useCallback
} from 'react';
import { Link as RouterLink } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import PerfectScrollbar from 'react-perfect-scrollbar';
import moment from 'moment';
import {
  Box,
  Card,
  CardHeader,
  Checkbox,
  Container,
  Collapse,
  Divider,
  InputAdornment,
  IconButton,
  FormControlLabel,
  Grid,
  Link,
  LinearProgress,
  TextField,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  TablePagination,
  SvgIcon,
  makeStyles
} from '@material-ui/core';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import {
  Search as SearchIcon
} from 'react-feather';
import _ from "lodash";
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import Label from 'src/components/Label';
import Header from './Header';

const sortOptions = [
  {
    value: 'age|desc',
    label: 'Age (newest first)'
  },
  {
    value: 'age|asc',
    label: 'Age (oldest first)'
  }
];

const descendingComparator = (a, b, orderBy) => {
  if (b[orderBy] < a[orderBy]) {
    return -1;
  }

  if (b[orderBy] > a[orderBy]) {
    return 1;
  }

  return 0;
};

const getComparator = (order, orderBy) => {
  return order === 'desc'
    ? (a, b) => descendingComparator(a, b, orderBy)
    : (a, b) => -descendingComparator(a, b, orderBy);
};

const applySort = (timelines, sort) => {
  const [orderBy, order] = sort.split('|');
  const comparator = getComparator(order, orderBy);
  const stabilizedThis = timelines.map((el, index) => [el, index]);

  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);

    if (order !== 0) return order;

    return a[1] - b[1];
  });

  return stabilizedThis.map((el) => el[0]);
};

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
                  <pre style={{ backgroundColor: "inherit", whiteSpace: "pre-wrap", overflowWrap: "anywhere" }} bgcolor="primary">
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
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
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

const LogListView = () => {
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
    console.log(q)
    getLogs(q, startDate, endDate, showOnlyErrors, page, limit)
  }, 300), []);

  const getLogs = async (q, startDate, endDate, showOnlyErrors, page, limit) => {
    setLoading(true);
    try {
      const payload = {
        "query": q,
        "timestamp_max": moment(startDate).format("yyyy-MM-DDTHH:mm:ssZ"),
        "timestamp_min": moment(endDate).format("yyyy-MM-DDTHH:mm:ssZ"),
        "only_errors": showOnlyErrors
      }
      const response = await axios.post(`/eywa/api/events/query?page=${page + 1}&per_page=${limit}`, payload)

      setLogs(response.data);
      setLoading(false);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get logs', {
        variant: 'error'
      });
    }
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
    getLogs(query, startDate, endDate, showOnlyErrors, page, limit);
  }, [page, limit, startDate, endDate, showOnlyErrors]);

  if (!logs && !loading) {
    return null;
  }

  return (
    <Page
      className={classes.root}
      title="Logs"
    >
      <Container maxWidth={false}>
        <Header />
        <Box mt={3}>
          <Card>
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
                <Grid item xs={12} sm={2}>
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
                <Grid item xs={12} sm={2}>
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
                <Grid item xs={6} sm={2}>
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
              <Fragment>
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
                            FUNCTIONS ID
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
              </Fragment>}
          </Card>
        </Box>
      </Container>
    </Page>
  );
};

export default LogListView;

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
  Card,
  CardHeader,
  InputAdornment,
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
import {
  Search as SearchIcon
} from 'react-feather';
import axios from 'src/utils/axios';
import Label from 'src/components/Label';
import useIsMountedRef from 'src/hooks/useIsMountedRef';

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

const useStyles = makeStyles(() => ({
  root: {},
  queryField: {
    width: 500
  }
}));

const Timelines = ({
  className,
  timelines,
  functionId,
  ...rest
}) => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const { enqueueSnackbar } = useSnackbar();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);
  const [filteredTimelines, setFilteredTimelines] = useState(timelines);
  const [loading, setLoading] = useState(true);

  const delayedQuery = useCallback(_.debounce((q) => {
    console.log(q)
    filterTimelines(q, functionId, page, limit)
  }, 300), []);

  const filterTimelines = async (query, functionId, page, limit) => {
    setLoading(true);
    try {
      let url = `/eywa/api/timeline?function_id=${functionId}&page=${page + 1}&per_page=${limit}`
      url = query !== '' ? url += `&query=${query}` : url;
      const response = await axios.get(url)

      setFilteredTimelines(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get timelines', {
        variant: 'error'
      });
    }
    setLoading(false);
  };

  const handleQueryChange = (event) => {
    event.persist();
    setQuery(event.target.value);
    delayedQuery(event.target.value);
  };

  const handleSortChange = (event) => {
    event.persist();
    setSort(event.target.value);
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage);
  };

  const handleLimitChange = (event) => {
    setLimit(parseInt(event.target.value));
  };

  useEffect(() => {
    if (isMountedRef) {
      filterTimelines(query, functionId, page, limit);
    }
  }, [page, limit]);

  const sortedTimelines = loading ? [] : applySort(filteredTimelines.objects, sort);

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Timeline Logs" />
      <Box
        p={2}
        minHeight={56}
        display="flex"
        alignItems="center"
      >
        <TextField
          className={classes.queryField}
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
          placeholder="Search timelines"
          value={query}
          variant="outlined"
        />
        <Box flexGrow={1} />
        <TextField
          label="Sort By"
          name="sort"
          onChange={handleSortChange}
          select
          SelectProps={{ native: true }}
          value={sort}
          variant="outlined"
        >
          {sortOptions.map((option) => (
            <option
              key={option.value}
              value={option.value}
            >
              {option.label}
            </option>
          ))}
        </TextField>
      </Box>
      {loading ? <LinearProgress />
        :
        <Fragment>
          <PerfectScrollbar>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    REQUEST ID
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
                {sortedTimelines.map((timeline) => (
                  <TableRow hover key={timeline.request_id} onClick={(event) => console.log(event)}>
                    <TableCell>
                      <Link
                        component={RouterLink}
                        to={`/app/timelines/${timeline.request_id}`}
                      >
                        <Typography
                          variant="h6"
                        // color="textPrimary"
                        >
                          {timeline.request_id}
                        </Typography>
                      </Link>
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
                      {moment(timeline.age).format('YYYY/MM/DD | hh:mm:ss')}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </PerfectScrollbar>
          <TablePagination
            component="div"
            count={filteredTimelines.total_count}
            onChangePage={handlePageChange}
            onChangeRowsPerPage={handleLimitChange}
            page={page}
            rowsPerPage={limit}
            rowsPerPageOptions={[5, 10, 25]}
          />

        </Fragment>
      }
    </Card>
  );
};

Timelines.propTypes = {
  className: PropTypes.string,
  timelines: PropTypes.object.isRequired,
  functionId: PropTypes.string.isRequired
};

export default Timelines;

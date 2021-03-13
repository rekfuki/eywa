import React, {
  Fragment,
  useState,
  useEffect
} from 'react';
import { Link as RouterLink } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import PerfectScrollbar from 'react-perfect-scrollbar';
import moment from 'moment';
import {
  Box,
  Card,
  Container,
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
import _ from "lodash";
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import Header from './Header';

const sortOptions = [
  {
    value: 'created_at|desc',
    label: 'Created at (newest first)'
  },
  {
    value: 'created_at|asc',
    label: 'Created at (oldest first)'
  },
  {
    value: 'updated_at|desc',
    label: 'Updated at (oldest first)'
  },
  {
    value: 'updated_at|asc',
    label: 'Updated at (oldest first)'
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

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3)
  },
  queryField: {
    width: 500
  }
}));

const applyFilters = (secrets, query) => {
  return secrets.filter((secret) => {
    let matches = true;

    if (query) {
      const properties = ['name', 'id', 'created_at', 'updated_at'];
      let containsQuery = false;

      properties.forEach((property) => {
        if (secret[property].toLowerCase().includes(query.toLowerCase())) {
          containsQuery = true;
        }
      });

      if (!containsQuery) {
        matches = false;
      }
    }

    return matches;
  });
};

const applyPagination = (secrets, page, limit) => {
  return secrets.slice(page * limit, page * limit + limit);
};

const SecretsListView = () => {
  const classes = useStyles();
  const { enqueueSnackbar } = useSnackbar();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);
  const [secrets, setSecrets] = useState({});
  const [loading, setLoading] = useState(true);

  const getSecrets = async () => {
    setLoading(true);
    try {
      const url = `/eywa/api/secrets`
      const response = await axios.get(url)

      setSecrets(response.data);
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get secrets', {
        variant: 'error'
      });
    }
    setLoading(false);
  };

  useEffect(() => {
    getSecrets();
  }, [])

  const handleQueryChange = (event) => {
    event.persist();
    setQuery(event.target.value);
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

  const filteredSecrets = loading ? [] : applyFilters(secrets.objects, query);
  const sortedSecrets = applySort(filteredSecrets, sort);
  const paginatedSecrets = applyPagination(sortedSecrets, page, limit)

  return (
    <Page
      className={classes.root}
      title="Secrets"
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
                placeholder="Search secrets"
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
            {loading
              ? <LinearProgress />
              :
              <Fragment>
                <PerfectScrollbar>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>
                          SECRET ID
                        </TableCell>
                        <TableCell>
                          SECRET NAME
                        </TableCell>
                        <TableCell>
                          CREATED AT
                        </TableCell>
                        <TableCell>
                          UPDATED AT
                        </TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {paginatedSecrets.map((secret) => (
                        <TableRow key={secret.id}>
                          <TableCell>
                            <Link
                              component={RouterLink}
                              to={`/app/secrets/${secret.id}`}
                            >
                              <Typography
                                variant="h6"
                              >
                                {secret.id}
                              </Typography>
                            </Link>
                          </TableCell>
                          <TableCell>
                            <Typography
                              variant="h6"
                            >
                              {secret.name}
                            </Typography>
                          </TableCell>
                          <TableCell>
                            {moment(secret.created_at).format('YYYY/MM/DD | hh:mm:ss')}
                          </TableCell>
                          <TableCell>
                            {moment(secret.updated_at).format('YYYY/MM/DD | hh:mm:ss')}
                          </TableCell>
                        </TableRow>

                      ))}
                    </TableBody>
                  </Table>
                </PerfectScrollbar>
                <TablePagination
                  component="div"
                  count={secrets.total_count}
                  onChangePage={handlePageChange}
                  onChangeRowsPerPage={handleLimitChange}
                  page={page}
                  rowsPerPage={limit}
                  rowsPerPageOptions={[5, 10, 25]}
                />
              </Fragment>
            }
          </Card>
        </Box>
      </Container>
    </Page>
  );
};

export default SecretsListView;

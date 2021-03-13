import React, {
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
  Container,
  InputAdornment,
  IconButton,
  Link,
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
import LinearProgress from '@material-ui/core/LinearProgress';
import DeleteIcon from '@material-ui/icons/DeleteOutline';
import {
  Search as SearchIcon,
  List as ListIcon
} from 'react-feather';
import _ from "lodash";
import axios from 'src/utils/axios';
import Page from 'src/components/Page';
import Label from 'src/components/Label';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Header from './Header';
import DeleteModal from './DeleteModal'

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
    value: 'expires_at|desc',
    label: 'Expires at (newest first)'
  },
  {
    value: 'expires_at|asc',
    label: 'Expires at (oldest first)'
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

const applySort = (tokens, sort) => {
  const [orderBy, order] = sort.split('|');
  const comparator = getComparator(order, orderBy);
  const stabilizedThis = tokens.map((el, index) => [el, index]);

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

const TokenListView = () => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const { enqueueSnackbar } = useSnackbar();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);
  const [tokens, setFilteredTokens] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);
  const [tokenToDelete, setTokenToDelete] = useState({});

  const handleDeleteModalOpen = (token) => {
    setTokenToDelete(token)
    setDeleteModalOpen(true);
  };

  const handleDeleteModalClose = () => {
    setDeleteModalOpen(false);
  };

  const handleDeleteModalExecute = async () => {
    try {
      await axios.delete('/eywa/api/tokens/' + tokenToDelete.id);
      if (isMountedRef.current) {
        const removeIndex = tokens.objects.map(function (token) { return token.id; }).indexOf(tokenToDelete.id);
        tokens.objects.splice(removeIndex, 1);

        setFilteredTokens(tokens);

        setDeleteModalOpen(false);

        enqueueSnackbar('Token deleted', {
          variant: 'success'
        });
      }
    } catch (err) {
      let message = "Delete request failed"
      console.error(err);
      enqueueSnackbar(message, {
        variant: 'error'
      });
    }
  };

  const delayedQuery = useCallback(_.debounce((q) => {
    console.log(q)
    getTokens(q, page, limit)
  }, 300), []);

  const getTokens = useCallback(async (query, page, limit) => {
    setLoading(true)
    try {
      let url = `/eywa/api/tokens?page=${page + 1}&per_page=${limit}`
      url = query !== '' ? url += `&query=${query}` : url;
      const response = await axios.get(url)

      console.log(response.data)
      if (isMountedRef.current) {
        setFilteredTokens(response.data);
      }

    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get tokens', {
        variant: 'error'
      });
    }
    setLoading(false)
  }, [isMountedRef]);

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

  const refresh = () => {
    getTokens("", 0, 10);
  }

  useEffect(() => {
    getTokens(query, page, limit);
  }, [getTokens])

  useEffect(() => {
    if (!isMountedRef.current) {
      getTokens(query, page, limit);
    }
  }, [page, limit]);

  if (!tokens) {
    return null;
  }

  const sortedTokens = applySort(tokens.objects, sort);

  return (
    <Page
      className={classes.root}
      title="Tokens"
    >
      <Container maxWidth={false}>
        <Header refresh={refresh} />
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
                placeholder="Search tokens"
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
              <div>
                <PerfectScrollbar>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>
                          TOKEN ID
                        </TableCell>
                        <TableCell align="center">
                          TOKEN NAME
                        </TableCell>
                        <TableCell align="center">
                          CREATED AT
                        </TableCell>
                        <TableCell align="center">
                          EXPIRES AT
                        </TableCell>
                        <TableCell align="right">
                          Actions
                        </TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {sortedTokens.map((token) => (
                        <TableRow key={token.id}>
                          <TableCell>
                            <Typography>{token.id}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            <Typography>{token.name}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            {moment(token.created_at * 1000).format('YYYY/MM/DD | hh:mm:ss')}
                          </TableCell>
                          <TableCell align="center">
                            {!token.expires_at * 1000 ? 'NEVER' : moment(token.expires_at).format('YYYY/MM/DD | hh:mm:ss')}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton onClick={() => handleDeleteModalOpen(token)}>
                              <SvgIcon fontSize="default" style={{ color: "red" }}>
                                <DeleteIcon />
                              </SvgIcon>
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </PerfectScrollbar>
                <TablePagination
                  component="div"
                  count={tokens.total_count}
                  onChangePage={handlePageChange}
                  onChangeRowsPerPage={handleLimitChange}
                  page={page}
                  rowsPerPage={limit}
                  rowsPerPageOptions={[5, 10, 25]}
                />
                <DeleteModal
                  tokenName={tokenToDelete.name}
                  onDelete={handleDeleteModalExecute}
                  onClose={handleDeleteModalClose}
                  open={isDeleteModalOpen}
                />
              </div>}
          </Card>
        </Box>
      </Container>
    </Page>
  );
};

export default TokenListView;

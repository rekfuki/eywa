import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Box,
  Card,
  Divider,
  IconButton,
  InputAdornment,
  Link,
  SvgIcon,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  makeStyles
} from '@material-ui/core';
import {
  Edit as EditIcon,
  ArrowRight as ArrowRightIcon,
  Search as SearchIcon
} from 'react-feather';
import Label from 'src/components/Label';

const sortOptions = [
  {
    value: 'updated_at|desc',
    label: 'Last update (newest first)'
  },
  {
    value: 'updated_at|asc',
    label: 'Last update (oldest first)'
  },
  {
    value: 'created_at|desc',
    label: 'Created at (newest first)'
  },
  {
    value: 'created_at|asc',
    label: 'Created at (oldest first)'
  }
];

const applyFilters = (fns, query) => {
  return fns.filter((fn) => {
    let matches = true;
    if (query) {
      const properties = ['id', 'short_name', 'image_id', 'image_name'];
      let containsQuery = false;

      properties.forEach((property) => {
        if (fn[property].toLowerCase().includes(query.toLowerCase())) {
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

const applyPagination = (fns, page, limit) => {
  return fns.slice(page * limit, page * limit + limit);
};

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

const applySort = (fns, sort) => {
  const [orderBy, order] = sort.split('|');
  const comparator = getComparator(order, orderBy);
  const stabilizedThis = fns.map((el, index) => [el, index]);

  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);

    if (order !== 0) return order;

    return a[1] - b[1];
  });

  return stabilizedThis.map((el) => el[0]);
};

const useStyles = makeStyles((theme) => ({
  root: {},
  queryField: {
    width: 500
  },
  bulkOperations: {
    position: 'relative'
  },
  bulkActions: {
    paddingLeft: 4,
    paddingRight: 4,
    marginTop: 6,
    position: 'absolute',
    width: '100%',
    zIndex: 2,
    backgroundColor: theme.palette.background.default
  },
  bulkAction: {
    marginLeft: theme.spacing(2)
  },
  avatar: {
    height: 42,
    width: 42,
    marginRight: theme.spacing(1)
  }
}));

const Results = ({
  className,
  fns,
  ...rest
}) => {
  const classes = useStyles();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);

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

  const filteredFns = applyFilters(fns, query);
  const sortedFns = applySort(filteredFns, sort);
  const paginatedFns = applyPagination(sortedFns, page, limit);

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Divider />
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
          placeholder="Search functions"
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
      <PerfectScrollbar>
        <Box minWidth={700}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  FUNCTION ID
                </TableCell>
                <TableCell align="center">
                  STATE
                </TableCell>
                <TableCell align="center">
                  NAME
                </TableCell>
                <TableCell align="center">
                  IMAGE ID
                </TableCell>
                <TableCell align="center">
                  IMAGE NAME
                </TableCell>
                <TableCell align="center">
                  CREATED AT
                </TableCell>
                <TableCell align="center">
                  UPDATED AT
                </TableCell>
                <TableCell align="right">
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {paginatedFns.map((fn) => {

                return (
                  <TableRow
                    hover
                    key={fn.id}
                  >
                    <TableCell>
                      <Box
                        display="flex"
                        alignItems="center"
                      >
                        <div>
                          <Link
                            component={RouterLink}
                            to={`/app/functions/${fn.id}`}
                            variant="h6"
                          >
                            {fn.id}
                          </Link>
                        </div>
                      </Box>
                    </TableCell>
                    <TableCell align="center">
                      <Label color={fn.available ? 'success' : 'error'}>
                        {fn.available
                          ? 'Available'
                          : fn.deleted_at !== undefined
                            ? "Terminating"
                            : "Unavailable"}
                      </Label>
                    </TableCell>
                    <TableCell align="center">
                      {fn.short_name}
                    </TableCell>
                    <TableCell align="center">
                      <Link
                        component={RouterLink}
                        to={`/app/images/${fn.image_id}/buildlogs`}
                        variant="h6"
                      >
                        {fn.image_id}
                      </Link>
                    </TableCell>
                    <TableCell align="center">
                      {fn.image_name}
                    </TableCell>
                    <TableCell align="center">
                      {fn.created_at}
                    </TableCell>
                    <TableCell align="center">
                      {fn.updated_at}
                    </TableCell>
                    <TableCell align="right">
                      <IconButton
                        component={RouterLink}
                        to={`/app/functions/${fn.id}/edit`}
                      >
                        <SvgIcon fontSize="small">
                          <EditIcon />
                        </SvgIcon>
                      </IconButton>
                      <IconButton
                        component={RouterLink}
                        to={`/app/functions/${fn.id}`}
                      >
                        <SvgIcon fontSize="small">
                          <ArrowRightIcon />
                        </SvgIcon>
                      </IconButton>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </Box>
      </PerfectScrollbar>
      <TablePagination
        component="div"
        count={filteredFns.length}
        onChangePage={handlePageChange}
        onChangeRowsPerPage={handleLimitChange}
        page={page}
        rowsPerPage={limit}
        rowsPerPageOptions={[5, 10, 25]}
      />
    </Card>
  );
};

Results.propTypes = {
  className: PropTypes.string,
  fns: PropTypes.array.isRequired
};

Results.defaultProps = {
  fns: []
};

export default Results;

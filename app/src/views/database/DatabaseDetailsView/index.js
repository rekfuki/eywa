import React, {
  useCallback,
  useState,
  useEffect
} from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import {
  Box,
  Container,
  InputAdornment,
  Card,
  Grid,
  TextField,
  LinearProgress,
  SvgIcon,
  Divider,
  makeStyles
} from '@material-ui/core';
import {
  Search as SearchIcon,
  List as ListIcon
} from 'react-feather';
import Page from 'src/components/Page';
import axios from 'src/utils/axios';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import Header from './Header';
import DatabaseInfo from './DatabaseInfo';
import CollectionInfo from './CollectionInfo';
import DeleteModal from './DeleteModal'

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

const sortOptions = [
  {
    value: 'storage_size|desc',
    label: 'Storage Size (largest first)'
  },
  {
    value: 'storage_size|asc',
    label: 'Storage Size (smallest first)'
  },
  {
    value: 'average_object_size|desc',
    label: 'Average Object Size (largest first)'
  },
  {
    value: 'average_object_size|asc',
    label: 'Average Object Size (smallest first)'
  },
  {
    value: 'index_count|desc',
    label: 'Index Count (largest first)'
  },
  {
    value: 'index_count|asc',
    label: 'Index Count (smallest first)'
  },
  {
    value: 'total_size|desc',
    label: 'Total Size (largest first)'
  },
  {
    value: 'total_size|asc',
    label: 'Total Size (smallest first)'
  },
  {
    value: 'total_index_size|desc',
    label: 'Total Index Size (largest first)'
  },
  {
    value: 'total_index_size|asc',
    label: 'Total Index Size (smallest first)'
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

const applySort = (images, sort) => {
  const [orderBy, order] = sort.split('|');
  const comparator = getComparator(order, orderBy);
  const stabilizedThis = images.map((el, index) => [el, index]);

  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);

    if (order !== 0) return order;

    return a[1] - b[1];
  });

  return stabilizedThis.map((el) => el[0]);
};

const applyFilters = (collections, query) => {
  return collections.filter((collection) => {
    let matches = true;

    if (query) {
      const properties = ['namespace'];
      let containsQuery = false;

      properties.forEach((property) => {
        if (collection[property].toLowerCase().includes(query.toLowerCase())) {
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

const DatabaseDetailsView = () => {
  const classes = useStyles();
  const history = useHistory();
  const isMountedRef = useIsMountedRef();
  const [databaseStats, setDatabaseStats] = useState(null);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);
  const { enqueueSnackbar } = useSnackbar();
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [collectionToDelete, setCollectionToDelete] = useState(null);

  const handleDeleteModalExecute = async () => {
    try {
      const collectionName = collectionToDelete.namespace.substring(collectionToDelete.namespace.indexOf('.') + 1)
      await axios.delete('/eywa/api/database/collections/' + collectionName);
      if (isMountedRef.current) {
        const removeIndex = databaseStats.collections_info
          .map(function (collection) { return collection.namespace; })
          .indexOf(collectionToDelete.namespace);

        databaseStats.collections_info.splice(removeIndex, 1);

        setDatabaseStats({ ...databaseStats });

        setDeleteModalOpen(false);

        enqueueSnackbar('Collection deleted', {
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

  const handleDeleteModalOpen = (collection) => {
    setCollectionToDelete(collection)
    setDeleteModalOpen(true);
  };

  const handleDeleteModalClose = () => {
    setDeleteModalOpen(false);
  };

  const handleQueryChange = (event) => {
    event.persist();
    setQuery(event.target.value);
  };

  const handleSortChange = (event) => {
    event.persist();
    setSort(event.target.value);
  };

  const getDatabaseStats = useCallback(async () => {
    setLoading(true)
    try {
      const response = await axios.get('/eywa/api/database')

      if (isMountedRef.current) {
        setDatabaseStats(response.data);
      }
      setLoading(false)
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get database stats', {
        variant: 'error'
      });
    }
  }, [isMountedRef]);

  useEffect(() => {
    getDatabaseStats();
  }, [getDatabaseStats]);

  const filteredCollections = loading ? [] : applyFilters(databaseStats.collections_info, query);
  const sortedCollections = applySort(filteredCollections, sort);

  return (
    <Page
      className={classes.root}
      title="Database Management"
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
                placeholder="Search collections"
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
          </Card>
          {loading ? <LinearProgress />
            :
            <div>
              <Divider />
              <Box mt={2}>
                <Grid
                  container
                  spacing={3}
                >
                  <Grid
                    item
                    lg={4}
                    md={6}
                    xl={3}
                    xs={12}
                  >
                    <DatabaseInfo databaseStats={databaseStats} />
                  </Grid>
                  {sortedCollections.map(collection => (
                    <Grid
                      key={collection.namespace}
                      item
                      lg={4}
                      md={6}
                      xl={3}
                      xs={12}
                    >
                      <CollectionInfo collection={collection} onDelete={handleDeleteModalOpen} />
                    </Grid>
                  ))}
                </Grid>
              </Box>
            </div>
          }
        </Box>
        <DeleteModal
          collection={collectionToDelete}
          onDelete={handleDeleteModalExecute}
          onClose={handleDeleteModalClose}
          open={isDeleteModalOpen}
        />
      </Container>
    </Page>
  );
};

export default DatabaseDetailsView;

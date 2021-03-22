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

function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

const ImagesListView = () => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const { enqueueSnackbar } = useSnackbar();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [query, setQuery] = useState('');
  const [sort, setSort] = useState(sortOptions[0].value);
  const [images, setFilteredImages] = useState({});
  const [loading, setLoading] = useState(true);
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);
  const [imageToDelete, setImageToDelete] = useState({});

  const handleDeleteModalOpen = (image) => {
    setImageToDelete(image)
    setDeleteModalOpen(true);
  };

  const handleDeleteModalClose = () => {
    setDeleteModalOpen(false);
  };

  const handleDeleteModalExecute = async () => {
    try {
      const response = await axios.delete('/eywa/api/images/' + imageToDelete.id);
      if (response.status === 204) {
        if (isMountedRef.current) {
          const removeIndex = images.objects.map(function (image) { return image.id; }).indexOf(imageToDelete.id);
          images.objects.splice(removeIndex, 1);

          setFilteredImages(images);

          setDeleteModalOpen(false);

          enqueueSnackbar('Image deleted', {
            variant: 'success'
          });
        }
      }
    } catch (err) {
      let message = "Delete request failed"
      if (err.response.status === 400) {
        message = "Cannot delete image that is building"
      }
      console.error(err);
      enqueueSnackbar(message, {
        variant: 'error'
      });
    }
  };

  const delayedQuery = useCallback(_.debounce((q) => {
    console.log(q)
    getImages(q, page, limit)
  }, 300), []);

  const getImages = useCallback(async (query, page, limit) => {
    setLoading(true)
    try {
      let url = `/eywa/api/images?page=${page + 1}&per_page=${limit}`
      url = query !== '' ? url += `&query=${query}` : url;
      const response = await axios.get(url)

      if (isMountedRef.current) {
        setFilteredImages(response.data);
      }

      setLoading(false)
    } catch (err) {
      console.error(err);
      enqueueSnackbar('Failed to get images', {
        variant: 'error'
      });
    }
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

  useEffect(() => {
    getImages(query, page, limit);
  }, [getImages])

  useEffect(() => {
    if (!isMountedRef.current) {
      getImages(query, page, limit);
    }
  }, [page, limit]);

  let sortedImages = []
  if (!loading) {
    sortedImages = applySort(images.objects, sort);
  }

  return (
    <Page
      className={classes.root}
      title="Images"
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
                placeholder="Search images"
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
                          IMAGE ID
                    </TableCell>
                        <TableCell align="center">
                          STATE
                    </TableCell>
                        <TableCell align="center">
                          IMAGE NAME
                    </TableCell>
                        <TableCell align="center">
                          RUNTIME
                    </TableCell>
                        <TableCell align="center">
                          VERSION
                    </TableCell>
                        <TableCell align="center">
                          SIZE
                    </TableCell>
                        <TableCell align="center">
                          CREATED AT
                    </TableCell>
                        <TableCell align="right">
                          Actions
                    </TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {sortedImages.map((image) => (
                        <TableRow key={image.id}>
                          <TableCell>
                            <Link
                              component={RouterLink}
                              to={`/app/images/${image.id}/buildlogs`}
                              variant="h6"
                            >
                              {image.id}
                            </Link>
                          </TableCell>
                          <TableCell align="center">
                            <Label
                              color={
                                image.state === "building"
                                  ? 'warning'
                                  : image.state === "success"
                                    ? 'success'
                                    : 'error'
                              }
                            >
                              {"BUILD " + image.state}
                            </Label>
                          </TableCell>
                          <TableCell align="center">
                            <Typography>{image.name}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            <Typography>{image.runtime}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            <Typography>{image.version}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            <Typography>{formatBytes(image.size, 3)}</Typography>
                          </TableCell>
                          <TableCell align="center">
                            {moment(image.created_at).format('YYYY/MM/DD | HH:mm:ss')}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton
                              component={RouterLink}
                              to={`/app/images/${image.id}/buildlogs`}
                            >
                              <SvgIcon fontSize="default">
                                <ListIcon />
                              </SvgIcon>
                            </IconButton>
                            <IconButton onClick={() => handleDeleteModalOpen(image)}>
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
                  count={images.total_count}
                  onChangePage={handlePageChange}
                  onChangeRowsPerPage={handleLimitChange}
                  page={page}
                  rowsPerPage={limit}
                  rowsPerPageOptions={[5, 10, 25]}
                />
                <DeleteModal
                  imageName={imageToDelete.name}
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

export default ImagesListView;

import React, { Fragment } from 'react';
import clsx from 'clsx';
import PerfectScrollbar from 'react-perfect-scrollbar';
import {
  Card,
  CardHeader,
  Divider,
  Table,
  IconButton,
  TableBody,
  TableCell,
  TableRow,
  Typography,
  SvgIcon,
  makeStyles
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/DeleteOutline';

const useStyles = makeStyles((theme) => ({
  root: {
    maxHeight: '50vh',
    overflowY: 'auto'
  },
  fontWeightMedium: {
    fontWeight: theme.typography.fontWeightMedium
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

const CollectionInfo = ({
  collection,
  className,
  onDelete,
  ...rest
}) => {
  const classes = useStyles();

  if (!collection) {
    return
  }

  const collectionName = collection.namespace.substring(collection.namespace.indexOf('.') + 1)
  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader
        style={{ textAlign: 'center', paddingBottom: 0 }}
        subheader={`collection`}
      />
      <CardHeader
        style={{ paddingTop: 0 }}
        title={`${collectionName}`}
        action={
          <IconButton onClick={() => onDelete(collection)}>
            <SvgIcon fontSize="default" style={{ color: "red" }}>
              <DeleteIcon />
            </SvgIcon>
          </IconButton>
        }
      />
      <Divider />
      <PerfectScrollbar>
        <Table>
          <TableBody>
            <TableRow>
              <TableCell>
                Total Size
            </TableCell>
              <TableCell>
                <Typography
                  variant="body2"
                  color="textSecondary"
                >
                  {formatBytes(collection.total_size)}
                </Typography>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                Storage Size
            </TableCell>
              <TableCell>
                <Typography
                  variant="body2"
                  color="textSecondary"
                >
                  {formatBytes(collection.storage_size)}
                </Typography>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                Total Index Size
            </TableCell>
              <TableCell>
                <Typography
                  variant="body2"
                  color="textSecondary"
                >
                  {formatBytes(collection.total_index_size)}
                </Typography>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                Averga Object Size
            </TableCell>
              <TableCell>
                <Typography
                  variant="body2"
                  color="textSecondary"
                >
                  {formatBytes(collection.average_object_size)}
                </Typography>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                Index Count
            </TableCell>
              <TableCell>
                <Typography
                  variant="body2"
                  color="textSecondary"
                >
                  {collection.index_count}
                </Typography>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <Typography
                  variant="h5"
                  color="textPrimary"
                >
                  Index Name
              </Typography>
              </TableCell>
              <TableCell>
                <Typography
                  variant="h5"
                  color="textPrimary"
                >
                  Size
              </Typography>
              </TableCell>
            </TableRow>
            {Object.entries(collection.index_sizes).map(([index, size]) => (
              <Fragment key={index}>
                <TableRow key={index}>
                  <TableCell>
                    {index}
                  </TableCell>
                  <TableCell>
                    {formatBytes(size)}
                  </TableCell>
                </TableRow >
              </Fragment>
            ))}
          </TableBody>
        </Table>
      </PerfectScrollbar>
    </Card>
  );
};

export default CollectionInfo;

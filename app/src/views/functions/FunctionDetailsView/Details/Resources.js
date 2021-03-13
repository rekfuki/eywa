
import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Card,
  CardHeader,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableRow,
  Typography,
  makeStyles,
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {},
  fontWeightMedium: {
    fontWeight: theme.typography.fontWeightMedium
  }
}));

const Resources = ({
  resources = {},
  className,
  ...rest
}) => {
  const classes = useStyles();

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Resources" />
      <Divider />
      <Table className={classes.table} size="small">
        <TableBody>
          <TableRow>
            <TableCell rowSpan={2}>
              Request
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                CPU
              </Typography>
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                {resources.min_cpu}
              </Typography>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                Memory
              </Typography>
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                {resources.min_memory}
              </Typography>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell rowSpan={2}>
              Limit
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                CPU
              </Typography>
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                {resources.max_cpu}
              </Typography>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                Memory
              </Typography>
            </TableCell>
            <TableCell>
              <Typography
                variant="body2"
                color="textSecondary"
              >
                {resources.max_memory}
              </Typography>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  );
};

Resources.propTypes = {
  className: PropTypes.string,
  resources: PropTypes.object.isRequired
};

export default Resources;

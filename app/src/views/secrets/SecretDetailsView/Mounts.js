import React from 'react';
import PropTypes from 'prop-types';
import { Link as RouterLink } from 'react-router-dom';
import clsx from 'clsx';
import {
  Card,
  CardHeader,
  Divider,
  Table,
  TableHead,
  Link,
  TableBody,
  TableCell,
  TableRow,
  Typography,
  makeStyles
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  root: {},
  fontWeightMedium: {
    fontWeight: theme.typography.fontWeightMedium
  }
}));

const Mounts = ({
  mounts,
  className,
  ...rest
}) => {
  const classes = useStyles();

  if (!mounts){
    return null;
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Mounted Functions" />
      <Divider />
      <Table>
        {(mounts.length == 0 &&
          <TableRow>
            <TableCell>
              <Typography
                variant="body2"
              >
                No Functions use this secret
              </Typography>
            </TableCell>
          </TableRow>
        ) ||
          <TableBody>
            {mounts.map((object, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Link
                    component={RouterLink}
                    to={`/app/functions/${object.id}`}
                  >
                    <Typography
                      variant="h6"
                    >
                      {object.id}
                    </Typography>
                  </Link>
                </TableCell>
                <TableCell>
                  {object.name}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        }
      </Table>
    </Card>
  );
};

export default Mounts;

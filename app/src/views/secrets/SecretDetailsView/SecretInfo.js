import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import moment from 'moment';
import {
  Card,
  CardHeader,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableRow,
  makeStyles
} from '@material-ui/core';


const useStyles = makeStyles(() => ({
  root: {}
}));

const SecretInfo = ({ className, secret, ...rest }) => {
  const classes = useStyles();

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Secret Info" />
      <Divider />
      <Table>
        <TableBody>
          <TableRow>
            <TableCell>
              ID
            </TableCell>
            <TableCell>
              {secret.id}
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              Name
            </TableCell>
            <TableCell>
              <div>{secret.name}</div>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              Mount Path
            </TableCell>
            <TableCell>
              <div>{"/var/faas/secrets/" + secret.name}</div>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              Total Mounts
            </TableCell>
            <TableCell>
              {secret.mounted_functions.length}
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              Updated at
            </TableCell>
            <TableCell>
              {moment(secret.updated_at).format('DD/MM/YYYY HH:MM')}
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              Created at
            </TableCell>
            <TableCell>
              {moment(secret.created_at).format('DD/MM/YYYY HH:MM')}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  );
};

SecretInfo.propTypes = {
  className: PropTypes.string
};

export default SecretInfo;

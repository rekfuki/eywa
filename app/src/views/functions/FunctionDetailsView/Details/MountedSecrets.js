import React from 'react';
import PropTypes from 'prop-types';
import { Link as RouterLink } from 'react-router-dom';
import clsx from 'clsx';
import {
  Card,
  CardHeader,
  Divider,
  Table,
  Link,
  TableHead,
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

const MountedSecrets = ({
  secrets,
  className,
  ...rest
}) => {
  const classes = useStyles();

  const secretsEmpty = (secrets === null || secrets.length === 0);
  console.log(secretsEmpty)
  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Mounted Secrets" />
      <Divider />
      <Table>
        {secretsEmpty ?
          <TableBody>
            <TableRow>
              <TableCell>
                <Typography
                  variant="body2"
                >
                  No secrets attached
                </Typography>
              </TableCell>
            </TableRow>
          </TableBody>
          :
          <TableBody>
            {secrets.map(s => (
              <TableRow key={s.id}>
                <TableCell>
                  <Link
                    component={RouterLink}
                    to={`/app/secrets/${s.id}`}
                  >
                    <Typography
                      variant="h6"
                    >
                      {s.id}
                    </Typography>
                  </Link>
                </TableCell>
                <TableCell>
                  <Typography
                    variant="body2"
                  >
                    {s.name}
                  </Typography>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        }
      </Table>
    </Card>
  );
};

MountedSecrets.propTypes = {
  className: PropTypes.string,
  secrets: PropTypes.array.isRequired
};

export default MountedSecrets;

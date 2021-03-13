import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Card,
  CardHeader,
  Divider,
  Table,
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

const EnvVars = ({
  envVars,
  className,
  ...rest
}) => {
  const classes = useStyles();

  const envEmpty = envVars === null || envVars.length === 0;
  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardHeader title="Environment Variables" />
      <Divider />
      <Table>
        {(envEmpty &&
          <TableRow>
            <TableCell>
              <Typography
                variant="body2"
              >
                No environment variables set
              </Typography>
            </TableCell>
          </TableRow>
        ) ||
          <TableBody>
            {Object.entries(envVars).map(([k, v]) => (
              <TableRow key={k}>
                <TableCell>
                  {k}
                </TableCell>
                <TableCell>
                  <Typography
                    variant="body2"
                    color="textSecondary"
                  >
                    {v}
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

EnvVars.propTypes = {
  className: PropTypes.string,
  envVars: PropTypes.object.isRequired
};

export default EnvVars;

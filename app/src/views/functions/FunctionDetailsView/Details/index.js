import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { Grid, makeStyles } from '@material-ui/core';
import FunctionInfo from './FunctionInfo';
import EnvVars from './EnvVars';
import MountedSecrets from './MountedSecrets';
import Resources from './Resources';

const useStyles = makeStyles(() => ({
  root: {}
}));

const Details = ({
  fn,
  className,
  ...rest
}) => {
  const classes = useStyles();

  return (
    <Grid
      className={clsx(classes.root, className)}
      container
      spacing={3}
      {...rest}
    >
      <Grid
        item
        lg={4}
        md={6}
        xl={3}
        xs={12}
      >
        <FunctionInfo fn={fn} />
      </Grid>
      <Grid
        item
        lg={4}
        md={6}
        xl={3}
        xs={12}
      >
        <EnvVars envVars={fn.env_vars || {}} />
      </Grid>
      <Grid
        item
        lg={4}
        md={6}
        xl={3}
        xs={12}
      >
        <MountedSecrets secrets={fn.secrets || []} />
      </Grid>
      {/* <Grid
        item
        lg={4}
        md={6}
        xl={3}
        xs={12}
      >
        <Resources resources={fn.resources} />
      </Grid> */}
    </Grid>
  );
};

Details.propTypes = {
  className: PropTypes.string,
  fn: PropTypes.object.isRequired
};

export default Details;

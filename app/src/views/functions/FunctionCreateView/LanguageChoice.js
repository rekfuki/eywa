import React, { useState } from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Box,
  Grid,
  Paper,
  FormHelperText,
  FormControlLabel,
  Typography,
  Radio,
  Button,
  makeStyles
} from '@material-ui/core';

const typeOptions = [
  {
    value: 'go',
    title: 'Go',
  },
];

const useStyles = makeStyles((theme) => ({
  root: {},
  typeOption: {
    alignItems: 'flex-start',
    display: 'flex',
    marginBottom: theme.spacing(2),
    padding: theme.spacing(2)
  },
  stepButton: {
    '& + &': {
      marginLeft: theme.spacing(2)
    }
  }
}));

const LanguageChoice = ({
  className,
  onBack,
  onNext,
  ...rest
}) => {
  const classes = useStyles();
  const [type, setType] = useState(typeOptions[0].value);
  const [isSubmitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (newType) => {
    setType(newType);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    
    try {
      setSubmitting(true);

      // NOTE: Make API request

      if (onNext) {
        onNext();
      }
    } catch (err) {
      console.error(err);
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Typography
        variant="h3"
        color="textPrimary"
      >
        Please select the language
      </Typography>
      <Box mt={2}>
        <Typography
          variant="subtitle1"
          color="textSecondary"
        >
          Programming language in which the function handler is written.
        </Typography>
      </Box>
        {typeOptions.map((typeOption) => (
          <Paper
            className={classes.typeOption}
            elevation={type === typeOption.value ? 10 : 1}
            key={typeOption.value}
          >
            <FormControlLabel
              control={
                <Radio
                  checked={type === typeOption.value}
                  onClick={() => handleChange(typeOption.value)}
                />
              }
              label={typeOption.title}
            />
          </Paper>
        ))}
      {error && (
        <Box mt={2}>
          <FormHelperText error>
            {error}
          </FormHelperText>
        </Box>
      )}
      <Box
        mt={6}
        display="flex"
      >
        {onBack && (
          <Button
            onClick={onBack}
            size="large"
          >
            Previous
          </Button>
        )}
        <Box flexGrow={1} />
        <Button
          color="secondary"
          disabled={isSubmitting}
          type="submit"
          variant="contained"
          size="large"
        >
          Next
        </Button>
      </Box>
    </form>
  );
};

LanguageChoice.propTypes = {
  className: PropTypes.string,
  onNext: PropTypes.func,
  onBack: PropTypes.func
};

LanguageChoice.defaultProps = {
  onNext: () => {},
  onBack: () => {}
};

export default LanguageChoice;

import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { useSnackbar } from 'notistack';
import {
  Box,
  Breadcrumbs,
  Button,
  IconButton,
  Grid,
  Link,
  List,
  ListItem,
  SvgIcon,
  Typography,
  TextField,
  Tooltip,
  makeStyles
} from '@material-ui/core';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import DeleteIcon from '@material-ui/icons/DeleteOutline';
import {
  Edit as EditIcon,
  Clipboard as ClipboardIcon
} from 'react-feather';
import useIsMountedRef from 'src/hooks/useIsMountedRef';
import DeleteModal from './DeleteModal'

const useStyles = makeStyles((theme) => ({
  root: {},
  deleteAction: {
    color: theme.palette.common.white,
    backgroundColor: theme.palette.error.main,
    '&:hover': {
      backgroundColor: theme.palette.error.dark
    }
  }
}));

const Header = ({
  className,
  fn,
  onDeleteFn,
  ...rest }) => {
  const classes = useStyles();
  const isMountedRef = useIsMountedRef();
  const { enqueueSnackbar } = useSnackbar();
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);


  const handleDeleteModalOpen = () => {
    setDeleteModalOpen(true);
  };

  const handleDeleteModalClose = () => {
    setDeleteModalOpen(false);
  };

  const handleDeleteModalExecute = () => {
    setDeleteModalOpen(false);
    onDeleteFn()
  }

  const copySyncUrl = (e) => {
    copyUrl(syncUrl);
  };

  const copyAsyncUrl = (e) => {
    copyUrl(asyncUrl);
  };

  const copyUrl = (text) => {
    navigator.clipboard.writeText(text)
    enqueueSnackbar('Copied', {
      variant: 'info'
    })
  }

  const baseUrl = `${window.location.protocol}//${window.location.host}/eywa/api/functions`
  const syncUrl = `${baseUrl}/sync/${fn.id}/`
  const asyncUrl = `${baseUrl}/async/${fn.id}/`

  return (
    <Grid
      container
      spacing={3}
      justify="space-between"
      className={clsx(classes.root, className)}
      {...rest}
    >
      <Grid item>
        <Breadcrumbs
          separator={<NavigateNextIcon fontSize="small" />}
          aria-label="breadcrumb"
        >
          <Link
            variant="body1"
            color="inherit"
            to="/app"
            component={RouterLink}
          >
            Dashboard
          </Link>
          <Link
            variant="body1"
            color="inherit"
            to="/app/functions"
            component={RouterLink}
          >
            Functions
          </Link>
          <Typography
            variant="body1"
            color="textPrimary"
          >
            {fn.id}
          </Typography>
        </Breadcrumbs>
        <Typography
          variant="h3"
          color="textPrimary"
        >
          {fn.name}
        </Typography>
        <Box mt={2}>
          <List>
            <ListItem>
              <Grid container spacing={1} alignItems="center">
                <Grid item xs={11}>
                  <TextField
                    id="filled-read-only-input"
                    label="Sync URL"
                    defaultValue={syncUrl}
                    fullWidth
                    margin="dense"
                    InputProps={{
                      readOnly: true,
                    }}
                  />
                </Grid>
                <Grid item xs={1}>
                  <Tooltip title="Copy to clipboard">
                    <IconButton onClick={copySyncUrl}>
                      <ClipboardIcon />
                    </IconButton>
                  </Tooltip>
                </Grid>
              </Grid>
            </ListItem>
            <ListItem>
              <Grid container spacing={1} alignItems="center">
                <Grid item xs={11}>
                  <TextField
                    id="filled-read-only-input"
                    label="Async URL"
                    defaultValue={asyncUrl}
                    fullWidth
                    margin="dense"
                    InputProps={{
                      readOnly: true,
                    }}
                  />
                </Grid>
                <Grid item xs={1}>
                  <Tooltip title="Copy to clipboard">
                    <IconButton onClick={copyAsyncUrl}>
                      <ClipboardIcon />
                    </IconButton>
                  </Tooltip>
                </Grid>
              </Grid>
            </ListItem>
          </List>
        </Box>
      </Grid>
      <Grid item>
        <Grid container spacing={3} justify="center">
          <Grid item>
            <Button
              disabled={fn.deleted_at !== undefined}
              className={classes.deleteAction}
              variant="contained"
              onClick={handleDeleteModalOpen}
              startIcon={
                <SvgIcon fontSize="small">
                  <DeleteIcon />
                </SvgIcon>
              }
            >
              Delete
            </Button>
            <DeleteModal
              fnName={fn.short_name}
              onDelete={handleDeleteModalExecute}
              onClose={handleDeleteModalClose}
              open={isDeleteModalOpen}
            />
          </Grid>
          <Grid item>
            <Button
              disabled={fn.deleted_at !== undefined}
              color="secondary"
              variant="contained"
              component={RouterLink}
              to={`/app/functions/${fn.id}/edit`}
              startIcon={
                <SvgIcon fontSize="small">
                  <EditIcon />
                </SvgIcon>
              }
            >
              Edit
            </Button>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
};

Header.propTypes = {
  className: PropTypes.string,
  fn: PropTypes.object.isRequired,
  onDeleteFn: PropTypes.func.isRequired
};

export default Header;

import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import moment from 'moment';
import {
  Box,
  Card,
  Container,
  Divider,
  Grid, 
  Typography,
  Paper,
  makeStyles
} from '@material-ui/core';
import Timeline from '@material-ui/lab/Timeline';
import TimelineItem from '@material-ui/lab/TimelineItem';
import TimelineSeparator from '@material-ui/lab/TimelineSeparator';
import TimelineConnector from '@material-ui/lab/TimelineConnector';
import TimelineContent from '@material-ui/lab/TimelineContent';
import TimelineOppositeContent from '@material-ui/lab/TimelineOppositeContent';
import TimelineDot from '@material-ui/lab/TimelineDot';

const useStyles = makeStyles((theme) => ({
  root: {},
  error: {
    color: theme.palette.common.white,
    backgroundColor: theme.palette.error.main
  },
  ok: {
    color: theme.palette.common.white,
    backgroundColor: theme.palette.success.main
  }
}));

const Details = ({
  timeline,
  className,
  ...rest
}) => {
  const classes = useStyles();

  return (
      <Card>
        <Grid
          className={clsx(classes.root, className)}
          container
          spacing={0}
          {...rest}
        >
          <Grid
            item
            xs={12}
            sm
          >
            <Box p={2}>
              <Typography
                variant="h5"
                color="textPrimary"
              >
                Method
              </Typography>
              <Typography
                variant="h5"
                color="textSecondary"
              >
                {timeline.method}
              </Typography>
            </Box>   
          </Grid>
          <Grid
            item
            xs={12}
            sm
          >
            <Box p={2}>
              <Typography
                variant="h5"
                color="textPrimary"
              >
                Response
              </Typography>
              <Typography
                variant="h5"
                color="textSecondary"
              >
                {timeline.response}
              </Typography>
            </Box>   
          </Grid>
          <Grid
            item
            xs={12}
            sm
          >
            <Box p={2}>
              <Typography
                variant="h5"
                color="textPrimary"
              >
                Duration
              </Typography>
              <Typography
                variant="h5"
                color="textSecondary"
              >
                {
                moment.duration(timeline.duration).minutes() > 0
                ?
                `${moment.duration(timeline.duration).minutes()} min`
                :
                moment.duration(timeline.duration).seconds() > 0
                ?
                `${moment.duration(timeline.duration).seconds()} s`
                :
                `${moment.duration(timeline.duration).milliseconds()} ms`
                        
              }
              </Typography>
            </Box>   
          </Grid>
          <Grid
            item
            xs={12}
            sm
          >
            <Box p={2}>
              <Typography
                variant="h5"
                color="textPrimary"
              >
                Created at
              </Typography>
              <Typography
                variant="h5"
                color="textSecondary"
              >
                {moment(timeline.age).format('YYYY/MM/DD | hh:mm:ss')}
              </Typography>
            </Box>   
          </Grid>
          <Grid
            item
            xs={12}
            sm
          >
            <Box p={2}>
              <Typography
                variant="h5"
                color="textPrimary"
              >
                Request ID
              </Typography>
              <Typography
                variant="h5"
                color="textSecondary"
              >
                {timeline.request_id}
              </Typography>
            </Box>   
          </Grid>
          <Grid item xs={12}>
            <Divider/>
          </Grid>
          <Grid item xs={12}>
            <Container maxWidth="md">
              <Timeline align="alternate">
                {timeline.events.map((event) => (
                  <TimelineItem key={event.name}>
                    <TimelineOppositeContent>
                      <Typography variant="body2" color="textSecondary">
                        {
                          moment.duration(event.duration).minutes() > 0
                          ?
                          `${moment.duration(event.duration).minutes()} min`
                          :
                          moment.duration(event.duration).seconds() > 0
                          ?
                          `${moment.duration(event.duration).seconds()} s`
                          :
                          `${moment.duration(event.duration).milliseconds()} ms`
                        }
                      </Typography>
                    </TimelineOppositeContent>
                    <TimelineSeparator>
                      <TimelineDot className={(event.is_error? classes.error : classes.ok)} style={{borderRadius: "0px"}}>
                        <Typography>
                          {event.response}
                        </Typography>
                      </TimelineDot>
                      <TimelineConnector />
                    </TimelineSeparator>
                    <TimelineContent>
                      <Typography variant="h6" component="h1">
                        {event.name}
                      </Typography>
                      <Typography>
                          {moment(event.timestamp).format('YYYY/MM/DD | hh:mm:ss')}
                      </Typography>
                    </TimelineContent>
                  </TimelineItem>
                ))}
              </Timeline>
            </Container>
          </Grid>
        </Grid>
      </Card>
  );
};

Details.propTypes = {
  className: PropTypes.string,
  timeline: PropTypes.object.isRequired
};

export default Details;

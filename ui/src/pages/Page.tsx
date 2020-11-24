import React from 'react';
import clsx from 'clsx';
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Tasks from './Tasks';
import Typography from '@material-ui/core/Typography';
import Add from './AddDialog';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      flexGrow: 1,
    },
    menuButton: {
      marginRight: theme.spacing(2),
    },
    title: {
      flexGrow: 1,
      textAlign: 'left',
    },
    addButton: {
      background: 'black',
    },
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
  }),
);

export default function PersistentDrawerLeft() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            Task scheduler
          </Typography>
          <Add />
        </Toolbar>
      </AppBar>
      <br />
      <main className={clsx(classes.content)}>
        <Tasks />
      </main>
    </div>
  );
}

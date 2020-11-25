import React from 'react';
import Snackbar from '@material-ui/core/Snackbar';
import { makeStyles, Theme } from '@material-ui/core/styles';

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    width: '100%',
    '& > * + *': {
      marginTop: theme.spacing(2),
    },
  },
}));

interface Props {
  open: boolean;
  message: string;
  onClose: () => void;
}

export default function ErrorBar({ open, message, onClose }: Props) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Snackbar
        open={open}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        autoHideDuration={12000}
        onClose={onClose}
        message={<div style={{color: 'red', fontWeight: 600}}>
          ERROR <br />
          {message}
        </div>}
        color={'red'}
      />
    </div>
  );
}

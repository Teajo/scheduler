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
  title: string;
  color: string;
  message: string;
  onClose: () => void;
}

export default function Bar({ open, title, color, message, onClose }: Props) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Snackbar
        open={open}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        autoHideDuration={6000}
        onClose={onClose}
        message={
          <div>
            <div style={{ color: color, fontWeight: 600 }}>
              {title}
            </div>
            {message}
          </div>
        }
      />
    </div>
  );
}

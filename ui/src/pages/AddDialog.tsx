import React, { useEffect } from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import AddIcon from '@material-ui/icons/Add';
import DateInput from './DateInput';
import PubSelect from './PubSelect';

export default function FormDialog() {
  const [open, setOpen] = React.useState(false);
  const [date, setDate] = React.useState(new Date().toISOString());
  const [publishers, setPublishers] = React.useState([]);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <div>
      <Button color="inherit" onClick={handleClickOpen} startIcon={<AddIcon />}>
        Add task
      </Button>
      <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">Add task</DialogTitle>
        <DialogContent>
          <DateInput date={date} onChange={setDate} />
        </DialogContent>
        <DialogContent>
          <PubSelect />
        </DialogContent>
        {/* <DialogContent>
          Configuration
        </DialogContent> */}
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleClose} color="primary">
            Add
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

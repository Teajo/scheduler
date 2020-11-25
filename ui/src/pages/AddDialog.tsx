import React, { useEffect } from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import AddIcon from '@material-ui/icons/Add';
import DateInput from './DateInput';
import PubSelect from './PubSelect';
import axios from 'axios';
import ErrorBar from './ErrorBar';
import SuccessBar from './SuccessBar';

export default function FormDialog() {
  const [open, setOpen] = React.useState(false);
  const [date, setDate] = React.useState(new Date().toISOString());
  const [publishers, setPublishers] = React.useState<any[]>([]);
  const [openErrorBar, setOpenErrorBar] = React.useState(false);
  const [openSuccessBar, setOpenSuccessBar] = React.useState(false);
  const [messageErrorBar, setMessageErrorBar] = React.useState('');
  const [messageSuccessBar, setMessageSuccessBar] = React.useState('');

  const handleClickOpen = () => {
    setPublishers([]);
    setOpen(true);
  };

  const handleClose = () => {
    setPublishers([]);
    setOpen(false);
  };

  const handleAdd = () => {
    axios.post(`http://127.0.0.1:3000/tasks/schedule`, {
      date,
      publishers,
    }, {headers: {
      'Content-Type': 'application/json'
    }})
    .then((res) => {
      setOpenSuccessBar(true);
      setMessageSuccessBar(res?.data?.message || 'Task created');
    })
    .catch((err) => {
      setOpenErrorBar(true);
      setMessageErrorBar(err?.response?.data?.error || 'Error during task creation');
    });
  };

  const handleAddPublisher = () => {
    setPublishers([ ...publishers, {
      name: '',
      settings: {},
      retryStrategy: {
        timeout: 25,
        exponential: true,
        limit: 5,
      }
    } ]);
  };

  const handleSettingsChange = (settings: any, index: number) => {
    const pubs = [...publishers];
    pubs[index] = settings;
    setPublishers(pubs);
  };

  return (
    <div>
      <ErrorBar open={openErrorBar} message={messageErrorBar} onClose={() => setOpenErrorBar(false)} />
      <SuccessBar open={openSuccessBar} message={messageSuccessBar} onClose={() => setOpenSuccessBar(false)} />
      <Button color="inherit" onClick={handleClickOpen} startIcon={<AddIcon />}>
        Add task
      </Button>
      <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">Add task</DialogTitle>
        <DialogContent>
          <DateInput date={date} onChange={setDate} />
        </DialogContent>
        <DialogContent>
          {
            publishers.map((pub, index) => (
              <div style={{marginTop: '10px'}}>
                <PubSelect data={pub} onChange={(settings) => handleSettingsChange(settings, index)} />
              </div>
            ))
          }
          <br />
          <Button variant="contained" onClick={handleAddPublisher}>Add publisher</Button>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleAdd} color="primary">
            Add
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

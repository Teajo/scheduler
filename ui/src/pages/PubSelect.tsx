import React, { useEffect } from 'react';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import axios from 'axios';
import { Label } from '@material-ui/icons';
import { TextField } from '@material-ui/core';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    formControl: {
      // margin: theme.spacing(1),
      minWidth: 300,
    },
    selectEmpty: {
      // marginTop: theme.spacing(2),
    },
  }),
);

export default function PubSelect() {
  const classes = useStyles();
  const [publishers, setPublishers] = React.useState<any>({});
  const [selectedPublishers, setSelectedPublishers] = React.useState<string[]>([]);

  useEffect(() => {
    axios.get(`http://127.0.0.1:3000/tasks/publishers`)
      .then((res) => {
        const { data } = res.data;
        setPublishers(data);
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);

  const handleChange = (ev: any) => {
    const pub = ev.target.value;
    setSelectedPublishers(pub);
  };

  return (
    <div>
      <FormControl className={classes.formControl}>
        <InputLabel id="demo-simple-select-label">Publisher</InputLabel>
        <Select
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={selectedPublishers}
          multiple
          onChange={handleChange}
        >
          {
            Object.keys(publishers).map(pub => (
              <MenuItem value={pub}>{pub}</MenuItem>
            ))
          }
        </Select>
      </FormControl>
      {
        selectedPublishers.map(pub => (
          <div>
            {
              Object.keys(publishers[pub]).map(field => (
                <div>
                  <FormControl className={classes.formControl}>
                    {
                      publishers[pub][field].Possible !== null &&
                      <>
                        <InputLabel id="test">{pub + ' ' + field}</InputLabel>
                        <Select
                          labelId="test"
                          id="test"
                          value={[]}
                          multiple
                        >
                          {
                            publishers[pub][field].Possible.map((pub: any) => (
                              <MenuItem value={pub}>{pub}</MenuItem>
                            ))
                          }
                        </Select>
                      </>
                    }
                    {
                      publishers[pub][field].Possible === null &&
                      <TextField id="standard-basic" label={pub + ' ' + field} />
                    }
                  </FormControl>
                </div>
              ))
            }
          </div>
        ))
      }
    </div>
  );
}

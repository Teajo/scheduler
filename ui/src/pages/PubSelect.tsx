import React, { useEffect } from 'react';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import { FormControlLabel, Switch } from '@material-ui/core';
import Select from '@material-ui/core/Select';
import axios from 'axios';
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

interface Props {
  id: string;
  label: string;
  field: any;
  value: any;
  onChange: (value: any) => void;
}

const FieldComponent = ({ id, label, field, value, onChange }: Props) => {
  switch(field.type) {
    case 'STRING':
      if (field.possible) {
        return (
          <>
            <InputLabel id="demo-simple-select-label">{label}</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={value || ''}
              onChange={(ev) => onChange(ev.target.value)}
            >
              {
                field.possible.map((pub: string) => (
                  <MenuItem value={pub}>{pub}</MenuItem>
                ))
              }
            </Select>
          </>
        );
      } else {
        return (
          <TextField 
            label={label} 
            value={value || ''} 
            type="text" 
            id={id}
            onChange={(ev) => onChange(ev.target.value)} 
          />
        );  
      }
    case 'JSON_STRING':
      return (
        <TextField 
          label={label} 
          value={value || ''} 
          type="text" 
          id={id}
          onChange={(ev) => onChange(ev.target.value)} 
        />
      );
    case 'INT':
      return (
        <TextField 
          label={label} 
          value={value || 0} 
          type="number" 
          id={id}
          onChange={(ev) => onChange(ev.target.value)}
        />
      );
    case 'BOOL':
      return (
        <FormControlLabel
          id={id}
          control={<Switch checked={value || false} onChange={(ev) => onChange(ev.target.checked)} />}
          label={label}
        />
      );
    default:
      return (
        <TextField 
          label={label} 
          value={value || ''} 
          type="text" 
          id={id} 
          onChange={(ev) => onChange(ev.target.value)}
        />
      );
  }
};

interface PubSelectProps {
  data: Record<string, any>;
  onChange: (settings: Record<string, any>) => void;
}

export default function PubSelect({ data, onChange }: PubSelectProps) {
  const classes = useStyles();
  const [publishers, setPublishers] = React.useState<any>({});

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
    const s = { ...data };
    s.publisher = pub;
    onChange(s);
  };

  return (
    <div>
      <FormControl className={classes.formControl}>
        <InputLabel id="demo-simple-select-label">Publisher</InputLabel>
        <Select
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={data.publisher}
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
        <div>
          {
            Object.keys(publishers[data.publisher] || []).map((field, index) => (
              <div>
                <FormControl className={classes.formControl} style={{marginTop: '5px'}}>
                  <FieldComponent 
                    id={'field'+index} 
                    label={`${data.publisher} ${field}`}
                    field={publishers[data.publisher][field]}
                    value={data.settings[field]}
                    onChange={(value) => {
                      const s = { ...data };
                      s.settings[field] = value;
                      onChange(s);
                    }}
                  />
                </FormControl>
              </div>
            ))
          }
        </div>
      }
    </div>
  );
}

import React from 'react';
import { DateTimePicker, MuiPickersUtilsProvider } from '@material-ui/pickers';
import DateFnsUtils from '@date-io/date-fns';

export interface Props {
  date: Date;
  label: string;
  onChange: (date: Date | null) => void;
}

export default function TimeFilter(props: Props) {
  return (
    <>
      <MuiPickersUtilsProvider utils={DateFnsUtils}>
        <DateTimePicker
          variant="inline"
          label={props.label}
          value={props.date}
          onChange={props.onChange}
          autoOk
          ampm={false}
        />
      </MuiPickersUtilsProvider>
    </>
  );
}

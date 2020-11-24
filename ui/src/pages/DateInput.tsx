import React, { useState } from 'react';
import MaskedInput from 'react-text-mask';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import FormControl from '@material-ui/core/FormControl';

function TextMaskCustom(props: any): any {
  const { inputRef, ...other } = props;

  return (
    <MaskedInput
      {...other}
      ref={(ref: any) => {
        inputRef(ref ? ref.inputElement : null);
      }}
      mask={[ /\d/, /\d/, /\d/, /\d/, '-', /\d/, /\d/, '-', /\d/, /\d/, 'T', /\d/, /\d/, ':', /\d/, /\d/, ':', /\d/, /\d/, '.', /\d/, /\d/, /\d/, 'Z'  ]}
      placeholderChar={'\u2000'}
      showMask
    />
  );
}

interface Props {
  date: string;
  onChange: (date: string) => void;
}

export default function DateInput(props: Props) {
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    props.onChange(event.target.value);
  };

  return (
    <div>
      <FormControl>
        <InputLabel htmlFor="formatted-text-mask-input">Date</InputLabel>
        <Input
          value={props.date}
          onChange={handleChange}
          name="textmask"
          id="formatted-text-mask-input"
          inputComponent={TextMaskCustom}
        />
      </FormControl>
    </div>
  );
}

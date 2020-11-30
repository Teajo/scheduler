import React from 'react';
import Bar from './Bar';

interface Props {
  open: boolean;
  message: string;
  onClose: () => void;
}

export default function SuccessBar({ open, message, onClose }: Props) {
  return (
    <Bar 
      open={open} 
      color={'green'} 
      title={'SUCCESS'} 
      onClose={onClose} 
      message={message} 
    />
  );
}

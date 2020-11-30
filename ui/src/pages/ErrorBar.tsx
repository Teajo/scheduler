import React from 'react';
import Bar from './Bar';

interface Props {
  open: boolean;
  message: string;
  onClose: () => void;
}

export default function ErrorBar({ open, message, onClose }: Props) {
  return (
    <Bar 
      open={open} 
      color={'red'} 
      title={'ERROR'} 
      onClose={onClose} 
      message={message} 
    />
  );
}

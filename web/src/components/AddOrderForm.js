import React, { useState } from 'react';
import { TextField, Button, MenuItem, Grid } from '@mui/material';

const AddOrderForm = ({ onClose, onCreate }) => {
  const [content, setContent] = useState('');
  const [status, setStatus] = useState('active');

  const handleContentChange = (e) => {
    setContent(e.target.value);
  };

  const handleStatusChange = (e) => {
    setStatus(e.target.value);
  };

  const handleCancel = () => {
    onClose();
  };

  const handleSubmit = () => {
    onCreate({ "content":content, "status":status });
    onClose();
  };

  return (
    <Grid container spacing={2}  >
      <Grid item xs={12}>
        <TextField
          label="Content"
          value={content}
          onChange={handleContentChange}
          fullWidth
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          select
          label="Status"
          value={status}
          onChange={handleStatusChange}
          fullWidth
        >
          <MenuItem value="active">Active</MenuItem>
          <MenuItem value="disabled">Disabled</MenuItem>
        </TextField>
      </Grid>
      <Grid item >
        
            
            <Button onClick={handleSubmit} color="primary" variant="contained">
            Add Order
            </Button>
            <Button onClick={handleCancel} color="secondary" variant="contained"> 
            Cancel
            </Button>
      </Grid>
    </Grid>
  );
};

export default AddOrderForm;

import React, { useState,useEffect } from 'react';
import { useSelector } from 'react-redux';
import AddOrderForm from './AddOrderForm.js';
import { useDispatch } from 'react-redux';
import axios from './api';
import config from '../config';

import {
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Typography,
  Grid,
  Button,
  TextField,
  MenuItem,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';

import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import moment from 'moment';

const OrderDialog = ({ patient, onClose }) => {
  const [orders, setOrders] = useState([]);
  const [isEditing, setIsEditing] = useState(false);
  const [editingOrder, setEditingOrder] = useState({ content: '', status: '' });
  const [expandedIndex, setExpandedIndex] = useState(null);
  const [editingIndex, setEditingIndex] = useState(null);
  const [isAddingOrder, setIsAddingOrder] = useState(false);

  const accessToken = useSelector(state => state.auth.accessToken);
  const refreshToken = useSelector(state => state.auth.refreshToken);

  const dispatch = useDispatch();
  const fetchOrders = async () => {
    try {   
      const response = await axios.get(`${config.apiUrl}/order-service/api/v1/patients/${patient.ID}/orders`, 
      {
        headers: {
          'Content-Type': 'application/json',
        }
      });
  
      if (!response.data || response.data.Code !== 0) {
          throw new Error('Failed to fetch patients');
      }
      
      
      if (!response.data || response.data.Code!=0 ) {
          throw new Error('Invalid response data');
      }
      setOrders(response.data.Data)
      return response.data.Data;
    } catch (error) {
        console.error('Error fetching patients:', error);
        setOrders([])
    }
  };

  useEffect(() => {
    fetchOrders(); 
  }, [patient.ID,accessToken]); 

  const formatDate = (dateString) => {
    return moment(dateString).format('YYYY-MM-DD HH:mm:ss');
  };

  const handleEditing = (index,content,status) => {
    if (content!=orders[index].content){
      editingOrder.content=content
    }

    if (status!=orders[index].status){
      editingOrder.status=status
    }

    setEditingOrder({...editingOrder});
    setIsEditing(true);
    setEditingIndex(index); 
  };

  const handleEditCancel = (index) => {
    setIsEditing(false);
    setEditingOrder({...orders[index]});
    setEditingIndex(null);
  };

  const handleEditSave = async () => {
    try {
      if (JSON.stringify(editingOrder) === JSON.stringify(orders[editingIndex])) {
        console.log('order not change');
        return;
      }

      const response = await axios.put(`${config.apiUrl}/order-service/api/v1/patients/${editingOrder.patient_id}/orders/${editingOrder.id}`, 
      editingOrder,{
        headers: {
          'Content-Type': 'application/json',
        }
      });
  
      if (!response.data || response.data.Code !== 0) {
          throw new Error('Failed to save changes');    
      }

      fetchOrders();
      setIsEditing(false);

      alert('Order successfully updated!');

    } catch (error) {
      console.error('Error saving changes:', error);
      alert('Failed to update order. Please try again later.');
    }
  };
  
  const handleAccordionClick = (index) => {
    setExpandedIndex(index === expandedIndex ? null : index);
    setEditingOrder({ ...orders[index]}); 
  };
  
  const handleCreateOrder = async (orderData) => {
    try {
      const response = await axios.post(`${config.apiUrl}/order-service/api/v1/patients/${patient.ID}/orders`, 
      orderData,{
        headers: {
          'Content-Type': 'application/json',
        }
      });
  
      if (!response.data || response.data.Code !== 0) {
        throw new Error('Failed to create order');
      }
  
      setIsAddingOrder(false);
      fetchOrders();
      alert('Order successfully created!');
    } catch (error) {
      console.error('Error creating order:', error);
      alert('Failed to create order. Please try again later.');
    }
  };

  
  return (
    <Dialog open onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Grid container justifyContent="space-between" alignItems="center">
          <Typography>{patient.Name}'s Orders</Typography>
          <Button color="primary" onClick={() => setIsAddingOrder(true)}>
            <AddIcon />
            Add Order
          </Button>
        </Grid>  
      </DialogTitle>
      
      <DialogContent>
        {isAddingOrder && (
          <Accordion expanded>
            <AccordionSummary>
              <Typography>Add New Order</Typography>
            </AccordionSummary>
            <AccordionDetails justify="flex-end">
              <AddOrderForm onClose={() => setIsAddingOrder(false)} onCreate={handleCreateOrder} />
            </AccordionDetails>
          </Accordion>
        )}
        {orders && orders.map((order, index) => (
          <Accordion 
            key={index}
            expanded={index === expandedIndex}
            onChange={() => handleAccordionClick(index)}>
            <AccordionSummary sx={{ '&:hover': { backgroundColor: '#f0f0f0' } }}>
              <Typography>{order.content}</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={9}>
                <TextField
                  label="Content"
                  value={ editingOrder.content}
                  onChange={(e) => handleEditing(index, e.target.value, editingOrder.status )}
                  fullWidth
                />
              </Grid>
              <Grid item xs={3}>
                
              </Grid>
              
              <Grid item xs={12}>
                <TextField
                  select
                  label="Status"
                  value={editingOrder.status}
                  onChange={(e) => handleEditing(index, editingOrder.content, e.target.value )}
                  fullWidth
                >
                  <MenuItem value="active">Active</MenuItem>
                  <MenuItem value="disabled">Disabled</MenuItem>
                </TextField>
                </Grid>
                <Grid item xs={12}>
                  <Typography variant="body1">Created At: {formatDate(order.created_at)}</Typography>
                </Grid>
                <Grid item xs={12}>
                  <Typography variant="body1">Doctor Name: {order.doctor_name}</Typography>
                </Grid>

    
              <Grid item xs={12}>
                {isEditing &&  (
                  <>
                    <Button onClick={() => handleEditSave(index)} color="primary" variant="contained">
                      <EditIcon />
                      Save
                    </Button>
                    <Button onClick={() => handleEditCancel(index)} color="secondary" variant="contained">
                      <DeleteIcon />
                      Cancel
                    </Button>
                  </>
                )}
              </Grid>

              </Grid>
            </AccordionDetails>
          </Accordion>
        ))}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default OrderDialog;
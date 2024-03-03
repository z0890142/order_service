import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchPatients, selectPatient ,resetSelectedPatient} from '../store/patientsSlice';
import { logout } from '../store/authSlice'; 
import { useNavigate } from 'react-router-dom';

import { List, ListItem, Card, Typography,Button } from '@mui/material';
import OrderDialog from './OrderDialog';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';

const PatientList = () => {
  const dispatch = useDispatch();
  const patients = useSelector((state) => state.patients.patients);
  const selectedPatient = useSelector(state => state.patients.selectedPatient);
  const loading = useSelector(state => state.patients.loading);
  const error = useSelector(state => state.patients.error);
  const accessToken = useSelector(state => state.auth.accessToken);

  const navigate = useNavigate();


  useEffect(() => {

    
    dispatch(fetchPatients());
  }, [dispatch]);

  const handlePatientClick = (patient) => {
    dispatch(selectPatient(patient));
  };

  const handleLogout = () => {
    dispatch(logout());
    navigate('/');
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  console.log(patients)

  return (
    <div style={{ padding: '20px' }}>
      <h2 style={{ marginBottom: '20px' }}>Patients List</h2>
      <List component="div">
        {patients && patients.map((patient) => (
          <Card
            key={patient.ID}
            style={{ marginRight: '100px', marginLeft: '100px', marginBottom: '10px', cursor: 'pointer' }}
            onClick={() => handlePatientClick(patient)}
            sx={{':hover': {backgroundColor: '#f0f0f0'}}}
          >
            <ListItem component="div" style={{ display: 'flex', alignItems: 'center' }}>
              <AccountCircleIcon style={{ marginRight: '10px' }} />
              <div>
                <Typography variant="body1">Name: {patient.Name}</Typography>
                <Typography variant="body1">Age: {patient.Age}</Typography>
                <Typography variant="body1">Gender: {patient.Gender === 1 ? 'Male' : 'Female'}</Typography>
              </div>
            </ListItem>
          </Card>
        ))}
      </List>
      {selectedPatient && (
        <OrderDialog
          patient={selectedPatient}
          onClose={() => dispatch(resetSelectedPatient())}
        />
      )}
      <Button variant="contained" color="primary" onClick={handleLogout} >Logout</Button>
    </div>
  );
};

export default PatientList;

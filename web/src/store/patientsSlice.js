import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from '../components/api';
import config from '../config';

const initialState = {
  patients: [],
  selectedPatient: null,
  loading: false,
  error: null
};


export const fetchPatients = createAsyncThunk(
  'patients/fetchPatients',
  async (_, { getState }) => {
      const { auth } = getState();
      try {          
          const response = await axios.get(`${config.apiUrl}/order-service/api/v1/patients`, 
          {
              headers: {
                'Content-Type': 'application/json',
              }
          });
          
          
          if (!response.data || response.data.Code !== 0) {
            throw new Error('Invalid response data');
          }
          
          return response.data.Data;
      } catch (error) {
          console.error('Error fetching patients:', error);
          throw error;
      }
  }
);
  

const patientsSlice = createSlice({
  name: 'patients',
  initialState,
  reducers: {
    selectPatient(state, action) {
      state.selectedPatient = action.payload;
    },
    resetSelectedPatient(state, action) {
      state.selectedPatient = null;
    },
    addOrder(state, action) {
      const { patientId, order } = action.payload;
      const patient = state.patients.find((p) => p.id === patientId);
      if (patient) {
        patient.orders.push(order);
      }
    },
    fetchPatientsSuccess(state, action) {
      state.loading = false;
      state.patients = action.payload;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
    .addCase(fetchPatients.pending, state => {
        state.loading = true;
        state.error = null;
        state.patients = null;
    })
    .addCase(fetchPatients.fulfilled, (state, action) => {
        state.loading = false;
        state.patients = action.payload;
        console.log(action.payload)
    })
    .addCase(fetchPatients.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message;
        state.patients = null;
    });
  },
});

export const { selectPatient, addOrder,resetSelectedPatient } = patientsSlice.actions;
export default patientsSlice.reducer;
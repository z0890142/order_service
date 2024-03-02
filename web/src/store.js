import { configureStore } from '@reduxjs/toolkit';
import patientsReducer from './store/patientsSlice';
import authReducer from './store/authSlice';
export default configureStore({
  reducer: {
    patients: patientsReducer,
    auth:authReducer
  }
});
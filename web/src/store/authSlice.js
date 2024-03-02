import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  accessToken: localStorage.getItem('accessToken') || null,
  refreshToken: localStorage.getItem('refreshToken') || null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setTokens(state, action) {
      const { accessToken, refreshToken } = action.payload;
      state.accessToken = accessToken;
      state.refreshToken=refreshToken;

      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
    },
    logout(state){
      state.accessToken = null;
      state.refreshToken = null;
      localStorage.clear();
    },
    setAccessToken(state,action){
      state.accessToken = action.payload;
      localStorage.setItem('accessToken', action.payload);
    }
    
    
  },
});

export const { setTokens,logout,setAccessToken } = authSlice.actions;
export default authSlice.reducer;

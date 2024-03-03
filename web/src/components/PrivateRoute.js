import React from 'react';
import { Navigate, Route } from 'react-router-dom';
import { useSelector } from 'react-redux';

const PrivateRoute = ({ Component }) => {
    const accessToken = localStorage.getItem('accessToken');
    const isAuthenticated = accessToken !== null;

    return isAuthenticated ? <Component /> : <Navigate to="/" />;
};

export default PrivateRoute;
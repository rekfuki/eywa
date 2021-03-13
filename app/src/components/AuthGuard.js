import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import PropTypes from 'prop-types';
import useAuth from '../hooks/useAuth';

const AuthGuard = ({ children }) => {
  const location = useLocation();
  const { isAuthenticated, checkAuth } = useAuth();

  if (!isAuthenticated) {
    window.location = "/login"
    return <div />
  }

  return (<>{children}</>);
};

AuthGuard.propTypes = {
  children: PropTypes.node
};

export default AuthGuard;

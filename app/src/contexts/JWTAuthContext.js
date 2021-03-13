import React, {
  createContext,
  useEffect,
  useReducer
} from 'react';
import SplashScreen from 'src/components/SplashScreen';
import axios from 'src/utils/axios';

const initialAuthState = {
  isAuthenticated: false,
  isInitialised: false,
  user: null
};

const reducer = (state, action) => {
  switch (action.type) {
    case 'INITIALISE': {
      const { isAuthenticated, user } = action.payload;

      return {
        ...state,
        isAuthenticated,
        isInitialised: true,
        user
      };
    }
    case 'LOGOUT': {
      return {
        ...state,
        isAuthenticated: false,
        user: null
      };
    }
    default: {
      return { ...state };
    }
  }
};

const AuthContext = createContext({
  ...initialAuthState,
  method: 'JWT',
  checkAuth: Promise.resolve(),
  logout: () => { }
});

export const AuthProvider = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, initialAuthState);

  const logout = async () => {
    try {
      await axios.post('/logout');
    } catch (err) {
      console.error(err);
    }
    dispatch({
      type: 'INITIALISE',
      payload: {
        isAuthenticated: false,
        user: null
      }
    });
  };

  const checkAuth = async () => {
    try {
      const response = await axios.get('/users/me');
      dispatch({
        type: 'INITIALISE',
        payload: {
          isAuthenticated: true,
          user: response.data
        }
      });
    } catch (err) {
      console.error(err);
      dispatch({
        type: 'INITIALISE',
        payload: {
          isAuthenticated: false,
          user: null
        }
      });
    }
  };

  useEffect(() => {
    checkAuth()
  }, []);

  if (!state.isInitialised) {
    return <SplashScreen />;
  }

  return (
    <AuthContext.Provider
      value={{
        ...state,
        method: 'JWT',
        checkAuth: checkAuth,
        logout: logout
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { createBrowserHistory } from 'history';
import { create } from 'jss';
import rtl from 'jss-rtl';
import MomentUtils from '@date-io/moment';
import { SnackbarProvider } from 'notistack';
import {
  jssPreset,
  StylesProvider,
  ThemeProvider
} from '@material-ui/core';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import GlobalStyles from 'src/components/GlobalStyles';
import ScrollReset from 'src/components/ScrollReset';
import { AuthProvider } from 'src/contexts/JWTAuthContext';
import useSettings from 'src/hooks/useSettings';
import { createTheme } from 'src/theme';
import routes, { renderRoutes } from 'src/routes';

const jss = create({ plugins: [...jssPreset().plugins, rtl()] });
const history = createBrowserHistory();

class DebugRouter extends BrowserRouter {
  constructor(props) {
    super(props);
    console.log('initial history is: ', JSON.stringify(this.history, null, 2))
    this.history.listen((location, action) => {
      console.log(
        `The current URL is ${location.pathname}${location.search}${location.hash}`
      )
      console.log(`The last navigation action was ${action}`, JSON.stringify(this.history, null, 2));
    });
  }
}

const App = () => {
  const { settings } = useSettings();

  const theme = createTheme({
    direction: settings.direction,
    responsiveFontSizes: settings.responsiveFontSizes,
    theme: settings.theme
  });

  return (
    <ThemeProvider theme={theme}>
      <StylesProvider jss={jss}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          <SnackbarProvider
            dense
            maxSnack={3}
          >
            <BrowserRouter history={history}>
              <AuthProvider>
                <GlobalStyles />
                <ScrollReset />
                {renderRoutes(routes)}
              </AuthProvider>
            </BrowserRouter>
          </SnackbarProvider>
        </MuiPickersUtilsProvider>
      </StylesProvider>
    </ThemeProvider>
  );
};

export default App;

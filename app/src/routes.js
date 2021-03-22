import React, {
  Suspense,
  Fragment,
  lazy
} from 'react';
import {
  Switch,
  Redirect,
  Route
} from 'react-router-dom';
import DashboardLayout from 'src/layouts/DashboardLayout';
import DocsLayout from 'src/layouts/DocsLayout';
import LoadingScreen from 'src/components/LoadingScreen';
import AuthGuard from 'src/components/AuthGuard';

export const renderRoutes = (routes = []) => (
  <Suspense fallback={<LoadingScreen />}>
    <Switch>
      {routes.map((route, i) => {
        const Guard = route.guard || Fragment;
        const Layout = route.layout || Fragment;
        const Component = route.component;

        return (
          <Route
            key={i}
            path={route.path}
            exact={route.exact}
            render={(props) => (
              <Guard>
                <Layout>
                  {route.routes
                    ? renderRoutes(route.routes)
                    : <Component {...props} />}
                </Layout>
              </Guard>
            )}
          />
        );
      })}
    </Switch>
  </Suspense>
);

const routes = [
  {
    exact: true,
    path: '/404',
    component: lazy(() => import('src/views/errors/NotFoundView'))
  },
  {
    path: '/app',
    guard: AuthGuard,
    layout: DashboardLayout,
    routes: [
      // BEGIN DASHBOARD
      {
        exact: true,
        path: '/app',
        component: () => <Redirect to="/app/dashboard" />
      },
      {
        exact: true,
        path: '/app/dashboard',
        component: lazy(() => import('src/views/reports/MainDashboard'))
      },
      // END DASHBOARD
      // BEGIN FUNCTIONS
      {
        exact: true,
        path: '/app/functions',
        component: lazy(() => import('src/views/functions/FunctionListView'))
      },
      {
        exact: true,
        path: '/app/functions/create',
        component: lazy(() => import('src/views/functions/FunctionCreateView'))
      },
      {
        exact: true,
        path: '/app/functions/:functionId',
        component: lazy(() => import('src/views/functions/FunctionDetailsView'))
      },
      {
        exact: true,
        path: '/app/functions/:functionId/edit',
        component: lazy(() => import('src/views/functions/FunctionEditView'))
      },
      // END FUNCTIONS
      // BEGIN TIMELINES
      {
        exact: true,
        path: '/app/timelines',
        component: lazy(() => import('src/views/timelines/TimelinesListView'))
      },
      {
        exact: true,
        path: '/app/timelines/:requestId',
        component: lazy(() => import('src/views/timelines/TimelineDetailsView'))
      },
      // END TIMELINES
      // BEGIN IMAGES
      {
        exact: true,
        path: '/app/images',
        component: lazy(() => import('src/views/images/ImagesListView'))
      },
      {
        exact: true,
        path: '/app/images/:imageId/buildlogs',
        component: lazy(() => import('src/views/images/ImageDetailsView'))
      },
      {
        exact: true,
        path: '/app/images/create',
        component: lazy(() => import('src/views/images/ImageCreateView'))
      },
      // END IMAGES
      // BEGIN SECRETS
      {
        exact: true,
        path: '/app/secrets',
        component: lazy(() => import('src/views/secrets/SecretsListView'))
      },
      {
        exact: true,
        path: '/app/secrets/create',
        component: lazy(() => import('src/views/secrets/SecretCreateView'))
      },
      {
        exact: true,
        path: '/app/secrets/:secretId',
        component: lazy(() => import('src/views/secrets/SecretDetailsView'))
      },
      {
        exact: true,
        path: '/app/secrets/:secretId/edit',
        component: lazy(() => import('src/views/secrets/SecretEditView'))
      },
      // END SECRETS
      // BEGIN TOKENS
      {
        exact: true,
        path: '/app/tokens',
        component: lazy(() => import('src/views/tokens/TokensListView'))
      },
      // END TOKENS
      // BEGIN DATABASE
      {
        exact: true,
        path: '/app/database',
        component: lazy(() => import('src/views/database/DatabaseDetailsView'))
      },
      // END DATABASE
      // BEGIN LOGS
      {
        exact: true,
        path: '/app/logs',
        component: lazy(() => import('src/views/logs/LogsListView'))
      }
      // END LOGS
    ]
  },
  {
    path: '/docs',
    layout: DocsLayout,
    guard: AuthGuard,
    routes: [
      {
        exact: true,
        path: '/docs',
        component: () => <Redirect to="/docs/overview/welcome" />
      },
      {
        exact: true,
        path: '/docs/*',
        component: lazy(() => import('src/views/docs/Docs'))
      }
      // {
      //   exact: true,
      //   path: '/docs/welcome',
      //   component: lazy(() => import('src/views/docs/WelcomeView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/getting-started',
      //   component: lazy(() => import('src/views/docs/GettingStartedView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/environment-variables',
      //   component: lazy(() => import('src/views/docs/EnvironmentVariablesView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/deployment',
      //   component: lazy(() => import('src/views/docs/DeploymentView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/api-calls',
      //   component: lazy(() => import('src/views/docs/APICallsView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/analytics',
      //   component: lazy(() => import('src/views/docs/AnalyticsView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/authentication',
      //   component: lazy(() => import('src/views/docs/AuthenticationView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/routing',
      //   component: lazy(() => import('src/views/docs/RoutingView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/settings',
      //   component: lazy(() => import('src/views/docs/SettingsView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/state-management',
      //   component: lazy(() => import('src/views/docs/StateManagementView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/theming',
      //   component: lazy(() => import('src/views/docs/ThemingView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/support',
      //   component: lazy(() => import('src/views/docs/SupportView'))
      // },
      // {
      //   exact: true,
      //   path: '/docs/changelog',
      //   component: lazy(() => import('src/views/docs/ChangelogView'))
      // },
      // {
      //   component: () => <Redirect to="/404" />
      // }
    ]
  },
  {
    path: '*',
    routes: [
      {
        exact: true,
        path: '/',
        component: () => <Redirect to="/app/dashboard" />
      },
      {
        component: () => <Redirect to="/404" />
      }
    ]
  }
];

export default routes;

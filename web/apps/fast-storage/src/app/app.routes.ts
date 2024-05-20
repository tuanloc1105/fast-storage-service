import { Route } from '@angular/router';
import { authGuard, noAuthGuard } from './core/guard';

export const appRoutes: Route[] = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'app',
  },
  {
    path: 'app',
    loadComponent: () =>
      import('./views/layout/layout.component').then((m) => m.LayoutComponent),
    canActivate: [authGuard],
  },
  {
    path: 'app/initializing',
    loadComponent: () =>
      import('./shared/ui/gateway.component').then((m) => m.GatewayComponent),
  },
  {
    path: 'auth/login',
    loadComponent: () =>
      import('./views/login/login.component').then((m) => m.LoginComponent),
    canActivate: [noAuthGuard],
  },
  {
    path: 'auth/register',
    loadComponent: () =>
      import('./views/register/register.component').then(
        (m) => m.RegisterComponent
      ),
    canActivate: [noAuthGuard],
  },
];

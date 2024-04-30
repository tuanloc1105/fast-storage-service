import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'app',
  },
  {
    path: 'app',
    loadComponent: () =>
      import('./views/layout/layout.component').then((m) => m.LayoutComponent),
  },
];

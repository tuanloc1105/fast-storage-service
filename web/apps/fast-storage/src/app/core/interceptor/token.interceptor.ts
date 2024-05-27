import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { LocalStorageJwtService } from '@app/shared/services';
import { environment } from 'environments/environment';

export const tokenInterceptor: HttpInterceptorFn = (req, next) => {
  const is_production = environment.production;

  let token: string | null = null;
  inject(LocalStorageJwtService)
    .getAccessToken()
    .subscribe((t) => (token = t));

  if (token) {
    req = req.clone({
      url: `${is_production ? process.env['BACKEND_URL'] : '/api'}${req.url}`,
      setHeaders: {
        Authorization: `Bearer ${token}`,
        'ngrok-skip-browser-warning': 'pass',
      },
    });
  } else {
    req = req.clone({
      url: `${is_production ? process.env['BACKEND_URL'] : '/api'}${req.url}`,
      setHeaders: {
        'ngrok-skip-browser-warning': 'pass',
      },
    });
  }
  return next(req);
};

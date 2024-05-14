import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { LocalStorageJwtService } from '@app/shared/services';
import { AuthStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { MessageService } from 'primeng/api';
import { catchError, throwError } from 'rxjs';

export const errorHandlerInterceptor: HttpInterceptorFn = (req, next) => {
  const messageService = inject(MessageService);
  const router = inject(Router);
  const authStore = inject(AuthStore);

  let refresh_token: string | null = null;
  inject(LocalStorageJwtService)
    .getRefreshToken()
    .subscribe((t) => (refresh_token = t));

  return next(req).pipe(
    catchError((error) => {
      if (
        error instanceof HttpErrorResponse &&
        !req.url.includes('auth/login') &&
        error.status === 401
      ) {
        if (refresh_token) {
          patchState(authStore, { isRefreshing: true });
          router.navigateByUrl('/app/initializing');
        }
        messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error.errorMessage,
        });
      }
      return throwError(() => error);
    })
  );
};

import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { ErrorHandlerStore } from '@app/store';
import { MessageService } from 'primeng/api';
import { catchError, throwError } from 'rxjs';

export const errorHandlerInterceptor: HttpInterceptorFn = (req, next) => {
  const errorHandlerStore = inject(ErrorHandlerStore);
  const messageService = inject(MessageService);

  return next(req).pipe(
    catchError((error) => {
      if (error instanceof HttpErrorResponse) {
        switch (error.status) {
          case 401:
            errorHandlerStore.handleError401(error);
            break;
          case 404:
            errorHandlerStore.handleError404(error);
            break;
          default:
            messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: error.error.errorMessage,
            });
            break;
        }
      }
      return throwError(() => error);
    })
  );
};

import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { MessageService } from 'primeng/api';

export interface ErrorHandlerState {
  errorCode: number;
  errorMessage: string | undefined;
  trace: string | undefined;
}

export const errorHandlerInitialState: ErrorHandlerState = {
  errorMessage: undefined,
  trace: undefined,
  errorCode: -1,
};

export const ErrorHandlerStore = signalStore(
  { providedIn: 'root' },
  withState<ErrorHandlerState>(errorHandlerInitialState),
  withMethods(
    (
      store,
      router = inject(Router),
      messageService = inject(MessageService)
    ) => ({
      handleError401: (error: HttpErrorResponse) => {
        patchState(store, {
          errorCode: error.error.errorCode,
          errorMessage: error.error.errorMessage,
          trace: error.error.trace,
        });
        messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error.errorMessage,
        });
        router.navigate(['auth/login']);
      },
      handleError404: (error: HttpErrorResponse) => {
        patchState(store, {
          errorCode: error.error.errorCode,
          errorMessage: error.error.errorMessage,
          trace: error.error.trace,
        });
        messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error.errorMessage,
        });
      },
    })
  )
);

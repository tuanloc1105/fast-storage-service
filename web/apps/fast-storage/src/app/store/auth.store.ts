import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '@app/data-access';
import {
  GetNewTokenRequest,
  LoginRequest,
  LogoutRequest,
  RegisterRequest,
  UserInfoResponse,
} from '@app/shared/model';
import { LocalStorageJwtService } from '@app/shared/services';
import { tapResponse } from '@ngrx/operators';
import { patchState, signalStore, withMethods, withState } from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { MessageService } from 'primeng/api';
import { pipe, switchMap } from 'rxjs';
import { StorageStore } from './storage.store';

type AuthState = {
  isLoggedIn: boolean;
  user: UserInfoResponse | null;
  isLoading: boolean;
  tryRefreshingToken: boolean;
};

const initialState: AuthState = {
  isLoggedIn: false,
  user: null,
  isLoading: false,
  tryRefreshingToken: false,
};

export const AuthStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods(
    (
      store,
      router = inject(Router),
      authService = inject(AuthService),
      localStorageService = inject(LocalStorageJwtService),
      messageService = inject(MessageService),
      storageStore = inject(StorageStore)
    ) => ({
      getUserInfo: rxMethod<void>(
        pipe(
          switchMap(() => {
            patchState(store, { isLoading: true });
            return authService.getUserInfo().pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { user: res.response, isLoggedIn: true });
                  if (storageStore.currentPath()) {
                    router.navigate(['app'], {
                      queryParams: { path: storageStore.currentPath() },
                    });
                  } else {
                    router.navigate(['app']);
                  }
                },
                error: () =>
                  patchState(store, {
                    isLoggedIn: false,
                    tryRefreshingToken: true,
                  }),
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      refreshToken: rxMethod<GetNewTokenRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return authService.getNewToken(payload).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, {
                    tryRefreshingToken: false,
                    isLoggedIn: true,
                  });
                  localStorageService.setItem({
                    access_token: res.response.accessToken,
                    refresh_token: res.response.refreshToken,
                  });
                  if (storageStore.currentPath()) {
                    router.navigate(['app'], {
                      queryParams: { path: storageStore.currentPath() },
                    });
                  } else {
                    router.navigate(['app']);
                  }
                },
                error: () => {
                  patchState(store, {
                    tryRefreshingToken: false,
                    isLoggedIn: false,
                  });
                  localStorageService.removeItem();
                  router.navigate(['auth/login']);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      register: rxMethod<RegisterRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return authService.register(payload).pipe(
              tapResponse({
                next: () => {
                  messageService.add({
                    severity: 'success',
                    summary: 'You have successfully registered',
                    detail: 'Please check your email to verify your account',
                  });
                  router.navigate(['auth/login']);
                },
                error: () => {
                  messageService.add({
                    severity: 'error',
                    summary: 'Error!',
                    detail: 'Something went wrong',
                  });
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      login: rxMethod<LoginRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return authService.login(payload).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { isLoggedIn: true });
                  localStorageService.setItem({
                    access_token: res.response.accessToken,
                    refresh_token: res.response.refreshToken,
                  });
                  if (localStorageService.getIsFirstTime()?.is) {
                    messageService.add({
                      severity: 'success',
                      summary: 'Welcome back!',
                      detail: 'You have successfully logged in',
                    });
                  }
                  router.navigate(['app']);
                },
                error: () => patchState(store, { isLoggedIn: false }),
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      logout: rxMethod<LogoutRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return authService.logout(payload).pipe(
              tapResponse({
                next: () => {
                  patchState(store, { isLoggedIn: false });
                  localStorageService.removeItem();
                  router.navigate(['auth/login']);
                },
                error: () => {
                  patchState(store, { isLoggedIn: false });
                  localStorageService.removeItem();
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
    })
  )
);

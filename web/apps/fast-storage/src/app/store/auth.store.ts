import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '@app/data-access';
import { LoginRequest, LogoutRequest } from '@app/shared/model';
import { LocalStorageJwtService } from '@app/shared/services';
import { tapResponse } from '@ngrx/operators';
import { patchState, signalStore, withMethods, withState } from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { pipe, switchMap, lastValueFrom } from 'rxjs';

type AuthState = {
  isLoggedIn: boolean;
  user: any | null;
  isLoading: boolean;
  isRefreshing: boolean;
};

const initialState: AuthState = {
  isLoggedIn: false,
  user: null,
  isLoading: false,
  isRefreshing: false,
};

export const AuthStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods(
    (
      store,
      router = inject(Router),
      authService = inject(AuthService),
      localStorageService = inject(LocalStorageJwtService)
    ) => ({
      async getUserInfo() {
        patchState(store, { isLoading: true });
        try {
          const info = await lastValueFrom(authService.getUserInfo());
          patchState(store, { user: info, isLoggedIn: true });
          router.navigate(['app']);
        } catch {
          patchState(store, { isLoggedIn: false });
          return Promise.reject();
        } finally {
          patchState(store, { isLoading: false });
        }
      },
      async refreshToken() {
        const refreshToken = await lastValueFrom(
          localStorageService.getRefreshToken()
        );

        if (!refreshToken) {
          patchState(store, { isLoggedIn: false });
          localStorageService.removeItem();
          router.navigate(['auth/login']);
          return Promise.reject();
        }

        patchState(store, { isRefreshing: true });
        try {
          const info = await lastValueFrom(
            authService.getNewToken({ request: { refreshToken } })
          );
          patchState(store, { isRefreshing: false, isLoggedIn: true });
          localStorageService.setItem({
            access_token: info.response.accessToken,
            refresh_token: info.response.refreshToken,
          });
          router.navigate(['app']);
        } catch {
          patchState(store, { isLoggedIn: false });
          localStorageService.removeItem();
          router.navigate(['auth/login']);
          return Promise.reject();
        }
      },
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
                error: () => patchState(store, { isLoggedIn: false }),
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
    })
  )
);

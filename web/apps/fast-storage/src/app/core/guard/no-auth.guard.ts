import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthStore } from '@app/store';

export const noAuthGuard: CanActivateFn = () => {
  const router = inject(Router);
  const isLoggedIn = inject(AuthStore).isLoggedIn();
  if (isLoggedIn) {
    router.navigate(['/app']);
  }
  return !isLoggedIn;
};

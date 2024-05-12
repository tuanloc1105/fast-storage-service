import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthStore } from '@app/store';

export const authGuard: CanActivateFn = () => {
  const router = inject(Router);
  const isLoggedIn = inject(AuthStore).isLoggedIn();
  if (!isLoggedIn) {
    router.navigate(['/auth/login']);
  }
  return isLoggedIn;
};

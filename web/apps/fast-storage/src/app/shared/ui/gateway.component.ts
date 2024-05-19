import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  effect,
  inject,
} from '@angular/core';
import { Router } from '@angular/router';
import { AuthStore } from '@app/store';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { LocalStorageJwtService } from '../services';

@Component({
  selector: 'app-gateway',
  standalone: true,
  imports: [CommonModule, ProgressSpinnerModule],
  template: `<div class="flex h-screen">
    <div class="m-auto flex items-center flex-col">
      <p-progressSpinner ariaLabel="loading" />
      @if(authStore.tryRefreshingToken()) {
      <span>We're refreshing your credential</span>
      } @else {
      <span>Please wait while we checking your credential</span>
      }
    </div>
  </div>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GatewayComponent implements OnInit {
  public authStore = inject(AuthStore);
  private readonly router = inject(Router);
  private readonly localStorageService = inject(LocalStorageJwtService);

  constructor() {
    effect(
      () => {
        if (this.authStore.tryRefreshingToken()) {
          this.localStorageService
            .getRefreshToken()
            .subscribe((refreshToken) => {
              if (!refreshToken) {
                this.router.navigateByUrl('/auth/login');
              } else {
                this.authStore.refreshToken({ request: { refreshToken } });
              }
            });
        }
      },
      { allowSignalWrites: true }
    );
  }

  ngOnInit(): void {
    if (!this.authStore.tryRefreshingToken()) {
      this.authStore.getUserInfo();
    }
  }
}

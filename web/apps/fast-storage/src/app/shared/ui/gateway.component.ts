import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  inject,
} from '@angular/core';
import { Router } from '@angular/router';
import { AuthStore } from '@app/store';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { lastValueFrom } from 'rxjs';
import { LocalStorageJwtService } from '../services';

@Component({
  selector: 'app-gateway',
  standalone: true,
  imports: [CommonModule, ProgressSpinnerModule],
  template: `<div class="flex h-screen">
    <div class="m-auto flex items-center flex-col">
      <p-progressSpinner ariaLabel="loading" />
      @if(authStore.isRefreshing()) {
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
  #router = inject(Router);
  #localStorageService = inject(LocalStorageJwtService);

  ngOnInit(): void {
    if (this.authStore.isRefreshing()) {
      this.refreshToken();
    } else {
      this.getUserInfo();
    }
  }

  async getUserInfo() {
    const refreshToken = await lastValueFrom(
      this.#localStorageService.getRefreshToken()
    );
    await this.authStore.getUserInfo().catch(() => {
      if (refreshToken) {
        this.refreshToken();
      } else {
        this.#router.navigate(['auth/login']);
      }
    });
  }

  async refreshToken() {
    await this.authStore.refreshToken();
  }
}

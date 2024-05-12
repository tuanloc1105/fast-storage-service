import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthStore } from '@app/store';
import { ProgressSpinnerModule } from 'primeng/progressspinner';

@Component({
  selector: 'app-gateway',
  standalone: true,
  imports: [CommonModule, ProgressSpinnerModule],
  template: `<div class="flex h-screen">
    <div class="m-auto flex items-center flex-col">
      <p-progressSpinner ariaLabel="loading" />
      <span>Please wait while we checking your credential</span>
    </div>
  </div>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GatewayComponent {
  #authStore = inject(AuthStore);
  #router = inject(Router);

  constructor() {
    this.getUserInfo();
  }

  async getUserInfo() {
    await this.#authStore
      .getUserInfo()
      .then(() => {
        this.#router.navigate(['/app']);
      })
      .catch(() => {
        this.#router.navigate(['/auth/login']);
      });
  }
}

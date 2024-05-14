import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { LocalStorageJwtService } from '@app/shared/services';
import { AuthStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { ConfirmationService } from 'primeng/api';
import { AvatarModule } from 'primeng/avatar';
import { ButtonModule } from 'primeng/button';
import { lastValueFrom } from 'rxjs';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [AvatarModule, ButtonModule],
  template: `
    <div class="flex flex-col items-center justify-between h-full">
      <p-avatar
        label="V"
        size="large"
        [style]="{ 'background-color': '#2196F3', color: '#ffffff' }"
      ></p-avatar>
      <div class="flex flex-col gap-10">
        <p-button
          icon="pi pi-folder"
          [text]="true"
          size="large"
          severity="secondary"
        ></p-button>
        <p-button
          icon="pi pi-cog"
          [text]="true"
          size="large"
          severity="secondary"
        ></p-button>
      </div>
      <p-button
        icon="pi pi-sign-out"
        severity="danger"
        [text]="true"
        size="large"
        (click)="handleLogout($event)"
      ></p-button>
    </div>
  `,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  #authStore = inject(AuthStore);
  #router = inject(Router);
  #confirmationService = inject(ConfirmationService);
  #localStorageJwtService = inject(LocalStorageJwtService);

  public async handleLogout(event: Event) {
    const refreshToken = await lastValueFrom(
      this.#localStorageJwtService.getRefreshToken()
    );
    if (refreshToken) {
      this.#confirmationService.confirm({
        target: event.target as EventTarget,
        message: 'Do you want to log out?',
        header: 'Logout Confirmation',
        icon: 'pi pi-info-circle',
        acceptButtonStyleClass: 'p-button-danger p-button-text',
        rejectButtonStyleClass: 'p-button-text p-button-text',

        accept: () => {
          this.#authStore.logout({ request: { refreshToken } });
        },
      });
    } else {
      patchState(this.#authStore, { isLoggedIn: false });
      this.#router.navigate(['auth/login']);
    }
  }
}

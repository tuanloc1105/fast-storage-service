import { ChangeDetectionStrategy, Component } from '@angular/core';
import { AvatarModule } from 'primeng/avatar';
import { ButtonModule } from 'primeng/button';

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
      ></p-button>
    </div>
  `,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {}

import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  effect,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { PasswordModule } from 'primeng/password';
import { DialogModule } from 'primeng/dialog';
import { FormsModule } from '@angular/forms';
import { DynamicDialogConfig, DynamicDialogRef } from 'primeng/dynamicdialog';
import { StorageStore } from '@app/store';
import { ButtonModule } from 'primeng/button';
import { FolderProtectionRequest } from '@app/shared/model';

@Component({
  selector: 'app-lock-folder',
  standalone: true,
  imports: [
    CommonModule,
    PasswordModule,
    DialogModule,
    FormsModule,
    ButtonModule,
  ],
  template: `<div class="flex items-center gap-3 mb-3">
      <label for="folder" class="font-semibold w-6rem">Password</label>
      <p-password [(ngModel)]="password" [feedback]="false" />
    </div>
    <div class="flex justify-end gap-2" #footer>
      <p-button
        label="Cancel"
        severity="secondary"
        (click)="closeNewFolderDialog()"
      />
      <p-button
        [label]="config.data.unlockFolder ? 'Unlock' : 'Lock'"
        (click)="submit()"
        [disabled]="!password"
        [loading]="storageStore.isLoading()"
      />
    </div>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LockFolderComponent implements OnInit {
  public storageStore = inject(StorageStore);

  private ref = inject(DynamicDialogRef);
  public config = inject(DynamicDialogConfig);

  public password = '';

  constructor() {
    effect(() => {
      if (
        this.storageStore.hasFileProtected() ||
        !this.storageStore.folderRequirePassword()
      ) {
        this.ref.close();
      }
    });
  }

  ngOnInit(): void {
    console.log(this.config.data.unlockFolder);
  }

  public submit() {
    if (this.config.data.unlockFolder) {
      const currentPath = this.storageStore.currentPath();
      this.storageStore.getDetailsDirectory({
        path: currentPath,
        type: 'detailFolder',
        credential: this.password,
      });
    } else {
      const payload: FolderProtectionRequest = {
        request: {
          folder: this.storageStore.currentPath(),
          credential: this.password,
          credentialType: 'PASSWORD',
        },
      };
      this.storageStore.folderProtection(payload);
    }
  }

  public closeNewFolderDialog() {
    this.ref.close();
  }
}

import {
  ChangeDetectionStrategy,
  Component,
  effect,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { DialogModule } from 'primeng/dialog';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { FormsModule } from '@angular/forms';
import { StorageStore } from '@app/store';
import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { patchState } from '@ngrx/signals';

@Component({
  selector: 'app-new-folder',
  standalone: true,
  imports: [
    CommonModule,
    DialogModule,
    ButtonModule,
    InputTextModule,
    FormsModule,
  ],
  template: `
    <div class="flex items-center gap-3 mb-3">
      <label for="folder" class="font-semibold w-6rem">Folder name</label>
      <input
        pInputText
        id="folder"
        class="flex-auto"
        autocomplete="off"
        [(ngModel)]="newFolderName"
      />
    </div>
    <div class="flex justify-end gap-2" #footer>
      <p-button
        label="Cancel"
        severity="secondary"
        (click)="closeNewFolderDialog()"
      />
      <p-button
        label="Save"
        (click)="addNewFolder()"
        [disabled]="!newFolderName"
        [loading]="storageStore.isLoading()"
      />
    </div>
  `,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NewFolderComponent {
  public storageStore = inject(StorageStore);

  private ref = inject(DynamicDialogRef);

  public newFolderName = '';

  constructor() {
    effect(
      () => {
        if (this.storageStore.hasNewFolder()) {
          this.closeNewFolderDialog();
          patchState(this.storageStore, { hasNewFolder: false });
        }
      },
      { allowSignalWrites: true }
    );
  }

  public addNewFolder() {
    if (this.storageStore.currentPath() !== '') {
      this.storageStore.createFolder(
        this.storageStore.currentPath() + '/' + this.newFolderName
      );
    } else {
      this.storageStore.createFolder(this.newFolderName);
    }
  }

  public closeNewFolderDialog() {
    this.ref.close();
  }
}

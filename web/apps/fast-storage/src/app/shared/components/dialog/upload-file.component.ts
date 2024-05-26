import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  effect,
  inject,
} from '@angular/core';
import { StorageStore } from '@app/store';
import { BadgeModule } from 'primeng/badge';
import { ButtonModule } from 'primeng/button';
import { DynamicDialogRef } from 'primeng/dynamicdialog';
import {
  FileSelectEvent,
  FileUploadEvent,
  FileUploadHandlerEvent,
  FileUploadModule,
} from 'primeng/fileupload';
import { ProgressBarModule } from 'primeng/progressbar';

@Component({
  selector: 'app-upload-file',
  standalone: true,
  imports: [
    CommonModule,
    FileUploadModule,
    ButtonModule,
    BadgeModule,
    ProgressBarModule,
    ButtonModule,
  ],
  template: `<p-fileUpload
    (onUpload)="onUpload($event)"
    [multiple]="true"
    [customUpload]="true"
    (uploadHandler)="uploadHandler($event)"
    (onSelect)="onSelect($event)"
  >
    <ng-template
      pTemplate="header"
      let-chooseCallback="chooseCallback"
      let-clearCallback="clearCallback"
      let-uploadCallback="uploadCallback"
    >
      <p-button
        label="Choose"
        icon="pi pi-plus"
        iconPos="right"
        (click)="chooseCallback()"
      />
      <p-button
        label="Upload"
        icon="pi pi-cloud-upload"
        iconPos="right"
        (click)="uploadCallback()"
        [loading]="storageStore.isLoading()"
        [disabled]="uploadedFiles.length === 0 || storageStore.isLoading()"
      />
      <p-button
        label="Cancel"
        icon="pi pi-times"
        iconPos="right"
        (click)="clearCallback()"
        [disabled]="uploadedFiles.length === 0 || storageStore.isLoading()"
      />
    </ng-template>
    <ng-template pTemplate="content">
      <ul *ngIf="uploadedFiles.length">
        <li *ngFor="let file of uploadedFiles">
          {{ file.name }} - {{ file.size }} bytes
        </li>
      </ul>
    </ng-template>
  </p-fileUpload>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UploadFileComponent {
  public storageStore = inject(StorageStore);

  private ref = inject(DynamicDialogRef);

  public uploadedFiles: File[] = [];

  constructor() {
    effect(
      () => {
        if (this.storageStore.hasNewFile()) {
          this.ref.close(true);
        }
      },
      { allowSignalWrites: true }
    );
  }

  public onUpload(event: FileUploadEvent) {
    for (const file of event.files) {
      this.uploadedFiles.push(file);
    }
  }

  public onSelect(event: FileSelectEvent) {
    this.uploadedFiles = event.currentFiles;
  }

  public uploadHandler(event: FileUploadHandlerEvent) {
    this.storageStore.uploadFile({
      files: event.files,
      folderLocation: this.storageStore.currentPath(),
    });
  }
}

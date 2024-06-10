import { Component, OnInit, effect, inject } from '@angular/core';
import {
  FolderDetailComponent,
  FolderTreeComponent,
  SidebarComponent,
} from '@app/components';
import { LockFolderComponent } from '@app/shared/components';
import { StorageStore } from '@app/store';
import { DialogService } from 'primeng/dynamicdialog';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [SidebarComponent, FolderTreeComponent, FolderDetailComponent],
  template: ` <div class="flex flex-row h-screen">
    <div class="basis-[5%] px-3 py-7">
      <app-sidebar></app-sidebar>
    </div>
    <div class="basis-1/3 px-4 py-7">
      <app-folder-tree></app-folder-tree>
    </div>
    <div class="basis-full px-5 py-7">
      <app-folder-detail></app-folder-detail>
    </div>
  </div>`,
  styles: [],
})
export class LayoutComponent implements OnInit {
  private readonly storageStore = inject(StorageStore);
  private readonly dialogService = inject(DialogService);

  ngOnInit(): void {
    this.storageStore.getSystemStorageStatus();
  }

  constructor() {
    effect(() => {
      if (this.storageStore.folderRequirePassword()) {
        this.dialogService.open(LockFolderComponent, {
          header: 'Folder require password',
          data: {
            unlockFolder: true,
          },
        });
      }
    });
  }
}

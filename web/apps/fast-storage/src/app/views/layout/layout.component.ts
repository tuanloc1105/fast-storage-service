import { CommonModule } from '@angular/common';
import { Component, OnInit, effect, inject } from '@angular/core';
import {
  FolderDetailComponent,
  FolderTreeComponent,
  SidebarComponent,
} from '@app/components';
import { BannerComponent, LockFolderComponent } from '@app/shared/components';
import { RouteParamsService } from '@app/shared/services';
import { BannerStore, StorageStore } from '@app/store';
import { DialogService } from 'primeng/dynamicdialog';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [
    CommonModule,
    SidebarComponent,
    FolderTreeComponent,
    FolderDetailComponent,
    BannerComponent,
  ],
  template: `<app-banner *ngIf="bannerStore.isBannerVisible()"></app-banner>
    <div class="flex flex-row h-screen">
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
  public bannerStore = inject(BannerStore);

  private readonly storageStore = inject(StorageStore);
  private readonly dialogService = inject(DialogService);
  private readonly routeParamsService = inject(RouteParamsService);

  ngOnInit(): void {
    this.storageStore.getSystemStorageStatus();
  }

  constructor() {
    effect(
      () => {
        if (this.storageStore.folderRequirePassword()) {
          this.dialogService.open(LockFolderComponent, {
            header: 'Folder require password',
            data: {
              unlockFolder: true,
            },
          });
        }
        if (this.storageStore.currentPath()) {
          this.routeParamsService.setRouteParams({
            path: this.storageStore.currentPath(),
          });
        }
        this.showBanner();
      },
      { allowSignalWrites: true }
    );
  }

  private showBanner() {
    const offerBanner = sessionStorage.getItem('offerBanner');
    const storageUsed = this.storageStore.status()?.used;
    const storageTotal = this.storageStore.status()?.maximunSize;

    if (storageUsed && storageTotal) {
      const percentage = (storageUsed / storageTotal) * 100;

      if (percentage >= 90) {
        this.bannerStore.showBanner(
          'Upgrade your storage',
          'You are running out of storage space',
          'Upgrade now'
        );
      } else if (offerBanner === 'true') {
        this.bannerStore.showBanner(
          'Get 50% off on your first purchase',
          'Use code: FIRST50',
          'Upgrade now'
        );
      }
    }
  }
}

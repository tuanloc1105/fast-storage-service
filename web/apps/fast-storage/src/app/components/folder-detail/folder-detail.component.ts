import { CommonModule, JsonPipe } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  DestroyRef,
  OnInit,
  QueryList,
  ViewChildren,
  effect,
  inject,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { StorageService } from '@app/data-access';
import {
  FileContentComponent,
  LockFolderComponent,
  NewFolderComponent,
  SearchComponent,
  UploadFileComponent,
} from '@app/shared/components';
import { AppearDirective } from '@app/shared/directives';
import { Directory } from '@app/shared/model';
import { ImageSrcPipe } from '@app/shared/pipe';
import { LocalStorageJwtService } from '@app/shared/services';
import { AppStore, StorageStore } from '@app/store';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import { heroPencilSquare, heroScissors } from '@ng-icons/heroicons/outline';
import { patchState } from '@ngrx/signals';
import { environment } from 'environments/environment';
import {
  ConfirmationService,
  MenuItem,
  MenuItemCommandEvent,
  MessageService,
} from 'primeng/api';
import { BreadcrumbItemClickEvent, BreadcrumbModule } from 'primeng/breadcrumb';
import { ButtonModule } from 'primeng/button';
import { ContextMenuModule } from 'primeng/contextmenu';
import { DialogService, DynamicDialogModule } from 'primeng/dynamicdialog';
import { GalleriaModule } from 'primeng/galleria';
import { IconFieldModule } from 'primeng/iconfield';
import { Inplace, InplaceModule } from 'primeng/inplace';
import { InputIconModule } from 'primeng/inputicon';
import { InputTextModule } from 'primeng/inputtext';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { SpeedDialModule } from 'primeng/speeddial';
import { TableModule } from 'primeng/table';
import { TooltipModule } from 'primeng/tooltip';
import { lastValueFrom, tap } from 'rxjs';

@Component({
  selector: 'app-folder-detail',
  standalone: true,
  imports: [
    ButtonModule,
    JsonPipe,
    BreadcrumbModule,
    TableModule,
    SpeedDialModule,
    DynamicDialogModule,
    CommonModule,
    ImageSrcPipe,
    ContextMenuModule,
    IconFieldModule,
    InputIconModule,
    InputTextModule,
    NgIconComponent,
    TooltipModule,
    InplaceModule,
    AppearDirective,
    GalleriaModule,
    ProgressSpinnerModule,
  ],
  templateUrl: './folder-detail.component.html',
  styleUrl: './folder-detail.component.scss',
  providers: [provideIcons({ heroScissors, heroPencilSquare })],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderDetailComponent implements OnInit {
  public storageStore = inject(StorageStore);
  public appStore = inject(AppStore);

  private readonly dialogService = inject(DialogService);
  private readonly confirmationService = inject(ConfirmationService);
  private readonly localStorageJwtService = inject(LocalStorageJwtService);
  private readonly messageService = inject(MessageService);
  private readonly storageService = inject(StorageService);
  private readonly destroyRef = inject(DestroyRef);

  public home: MenuItem | undefined;
  public selectedDirectory: { directory: Directory; rowIndex: number } | null =
    null;
  public checkedDirectories: Directory[] = [];

  public tableContextMenu: MenuItem[] = [];
  public speedDialItems: MenuItem[] = [];

  public viewImages = false;
  public activeViewImageIndex = 0;
  public responsiveOptions: { breakpoint: string; numVisible: number }[] = [
    {
      breakpoint: '1500px',
      numVisible: 5,
    },
    {
      breakpoint: '1024px',
      numVisible: 3,
    },
    {
      breakpoint: '768px',
      numVisible: 2,
    },
    {
      breakpoint: '560px',
      numVisible: 1,
    },
  ];

  @ViewChildren(Inplace) inplaces!: QueryList<Inplace>;

  constructor() {
    effect(
      () => {
        if (
          this.storageStore.hasNewFolder() ||
          this.storageStore.hasNewFile() ||
          this.storageStore.hasFileRemoved() ||
          this.storageStore.hasFileRenamed()
        ) {
          this.storageStore.getDetailsDirectory({
            path: this.storageStore.currentPath(),
            type: 'detailFolder',
          });
        }

        if (this.storageStore.fileContent()) {
          const ref = this.dialogService.open(FileContentComponent, {
            header: this.selectedDirectory?.directory.name,
            data: {
              content: this.storageStore.fileContent(),
            },
            resizable: true,
            width: '50vw',
          });

          ref.onClose.subscribe(() => {
            patchState(this.storageStore, {
              fileContent: '',
            });
          });
        }
      },
      { allowSignalWrites: true }
    );
  }

  ngOnInit(): void {
    this.home = {
      icon: 'pi pi-home',
      command: () =>
        this.storageStore.getDetailsDirectory({
          path: '',
          type: 'detailFolder',
        }),
    };
    this.speedDialItems = [
      {
        icon: 'pi pi-refresh',
        command: () => console.log('Update'),
      },
      {
        icon: 'pi pi-folder-plus',
        command: () => this.addNewFolder(),
      },
      {
        icon: 'pi pi-upload',
        command: () => this.uploadFile(),
      },
      {
        icon: 'pi pi-lock',
        command: () => this.lockFolder(),
      },
    ];

    this.tableContextMenu = [
      {
        label: 'Change name',
        icon: 'pi pi-fw pi-eraser',
        command: () => {
          this.inplaces.forEach((inplace) => inplace.deactivate());
          this.inplaces
            .toArray()
            [this.selectedDirectory?.rowIndex ?? 0].activate();
        },
      },
      {
        label: 'Download',
        icon: 'pi pi-fw pi-download',
        command: () =>
          this.downloadFile(
            this.selectedDirectory && this.selectedDirectory?.directory
          ),
      },
      {
        label: 'Delete',
        icon: 'pi pi-fw pi-times',
        command: (e) => this.removeFile(e),
      },
    ];

    this.storageStore.getDetailsDirectory({ path: '', type: 'detailFolder' });
  }

  public handleBreadcrumb(event: BreadcrumbItemClickEvent): void {
    const item = event.item.label;
    const index = this.storageStore
      .breadcrumb()
      .findIndex((b) => b.label === item);
    patchState(this.storageStore, {
      breadcrumb: this.storageStore.breadcrumb().slice(0, index + 1),
    });
  }

  public confirmChangeName(newName: string, element: Inplace) {
    element.deactivate();
    const oldFolderLocationName =
      this.storageStore.currentPath() + this.selectedDirectory?.directory.name;
    const newFolderLocationName = this.storageStore.currentPath() + newName;

    this.storageStore.renameFileOrFolder({
      request: {
        oldFolderLocationName,
        newFolderLocationName,
      },
    });
  }

  public retrieveDirectory(directory: Directory, rowIndex: number): void {
    if (directory.type === 'folder') {
      const newPath = this.storageStore.currentPath() + '/' + directory.name;
      patchState(this.storageStore, {
        currentPath: newPath,
        breadcrumb: newPath
          .split('/')
          .filter((path) => path)
          .map((p) => ({
            label: p,
            command: () => {
              const nextPath = newPath
                .split('/')
                .slice(0, newPath.split('/').indexOf(p) + 1)
                .join('/');
              this.storageStore.getDetailsDirectory({
                path: nextPath,
                type: 'detailFolder',
              });
            },
            styleClass: 'cursor-pointer',
          })),
      });
      this.storageStore.getDetailsDirectory({
        path: newPath,
        type: 'detailFolder',
      });
    } else {
      this.selectedDirectory = { directory, rowIndex };
      if (directory.extension.toLowerCase().match(/(jpg|jpeg|png|gif)$/)) {
        const imageIndex = this.storageStore
          .imagesPool()
          .findIndex((image) => image.title === directory.name);
        this.activeViewImageIndex = imageIndex;
        if (!this.isImageReady(directory.name)) {
          this.retrieveImage(directory.name);
        } else {
          this.viewImages = true;
        }
      } else {
        this.storageStore.readFileContent({
          fileNameToRead: directory.name,
          locationToRead: this.storageStore.currentPath(),
        });
      }
    }
  }

  public downloadFiles(): void {
    this.checkedDirectories.forEach((directory) => {
      this.downloadFile(directory);
    });
  }

  public handleCopy() {
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: 'File(s) copied',
    });
  }

  public handleCut() {
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: 'File(s) cut',
    });
  }

  public handlePaste() {
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: 'File(s) pasted',
    });
  }

  public handleRename() {
    this.messageService.add({
      severity: 'success',
      summary: 'Success',
      detail: 'File(s) renamed',
    });
  }

  public deleteFiles(event: Event): void {
    this.confirmationService.confirm({
      target: event.target as EventTarget,
      message: 'Do you want to delete these file(s)?',
      header: 'Delete Confirmation',
      icon: 'pi pi-info-circle',
      acceptButtonStyleClass: 'p-button-danger p-button-text',
      rejectButtonStyleClass: 'p-button-text p-button-text',
    });
  }

  public handleSearch(): void {
    const dialogRef = this.dialogService.open(SearchComponent, {
      position: 'top',
      showHeader: false,
      width: '700px',
      contentStyle: { borderRadius: '12px', padding: '8px' },
      dismissableMask: true,
    });

    dialogRef.onClose.subscribe(() => {
      patchState(this.storageStore, {
        searchResults: [],
      });
    });
  }

  public onImageViewChange(index: number) {
    const imageName = this.storageStore.imagesPool()[index].title;
    if (!this.isImageReady(imageName)) {
      this.retrieveImage(imageName);
    }
  }

  private isImageReady(imageName: string): boolean {
    return !!this.storageStore
      .imagesPool()
      .find((image) => image.title === imageName)?.itemImageSrc;
  }

  private retrieveImage(imageFileName: string): void {
    patchState(this.storageStore, {
      isLoading: true,
    });
    this.storageService
      .showImage({
        request: {
          imageFileName,
          folderLocation: this.storageStore.currentPath(),
        },
      })
      .pipe(
        tap((res) => {
          patchState(this.storageStore, {
            imagesPool: this.storageStore.imagesPool().map((image) => {
              if (image.title === imageFileName) {
                return {
                  ...image,
                  itemImageSrc: res.response.data,
                };
              }
              return image;
            }),
            isLoading: false,
          });
          this.viewImages = true;
        }),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe();
  }

  private lockFolder(): void {
    this.dialogService.open(LockFolderComponent, {
      header: 'Lock folder',
    });
  }

  private addNewFolder(): void {
    this.dialogService.open(NewFolderComponent, {
      header: 'Create new folder',
    });
  }

  private uploadFile(): void {
    const uploadFileDialogRef = this.dialogService.open(UploadFileComponent, {
      header: 'Upload file',
      width: '50vw',
    });

    uploadFileDialogRef.onClose.subscribe((res) => {
      if (res) {
        this.storageStore.getDetailsDirectory({
          path: this.storageStore.currentPath(),
          type: 'detailFolder',
        });
      }
    });
  }

  private async downloadFile(directory: Directory | null) {
    const accessToken = await lastValueFrom(
      this.localStorageJwtService.getAccessToken()
    );
    if (!accessToken || !directory) return;

    if (directory.type === 'folder') {
      window.location.href = `${
        environment.apiUrl
      }/storage/download_folder?locationToDownload=${this.storageStore.currentPath()}&token=${accessToken}`;
      return;
    }

    window.open(
      `${environment.apiUrl}/storage/download_file?fileNameToDownload=${
        directory?.name
      }&locationToDownload=${this.storageStore.currentPath()}&token=${accessToken}`,
      '_blank'
    );
  }

  private removeFile(event: MenuItemCommandEvent): void {
    this.confirmationService.confirm({
      target: event.originalEvent?.target as EventTarget,
      message: 'Do you want to delete this?',
      header: 'Delete Confirmation',
      icon: 'pi pi-info-circle',
      acceptButtonStyleClass: 'p-button-danger p-button-text',
      rejectButtonStyleClass: 'p-button-text p-button-text',
      accept: () => {
        this.storageStore.removeFile({
          request: {
            fileNameToRemove: this.selectedDirectory?.directory.name ?? '',
            locationToRemove: this.storageStore.currentPath(),
          },
        });
      },
    });
  }
}

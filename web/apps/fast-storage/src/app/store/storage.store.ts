import { HttpErrorResponse } from '@angular/common/http';
import { computed, inject } from '@angular/core';
import { StorageService } from '@app/data-access';
import {
  COLOR_STATUS_STORAGE,
  DEFAULT_STATUS_STORAGE_COLOR,
} from '@app/shared/constant';
import {
  CheckFolderProtectionRequest,
  CutOrCopyRequest,
  Directory,
  DownloadFileRequest,
  FolderProtectionRequest,
  ImageViewer,
  ReadFileRequest,
  RemoveFileRequest,
  RenameRequest,
  SearchRequest,
  StorageStatus,
  UploadFileRequest,
} from '@app/shared/model';
import { tapResponse } from '@ngrx/operators';
import {
  patchState,
  signalStore,
  withComputed,
  withMethods,
  withState,
} from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { MenuItem, MessageService, TreeNode } from 'primeng/api';
import { pipe, switchMap } from 'rxjs';

type StorageState = {
  status: StorageStatus | null;
  allDirectorys: Directory[];
  subMenuDirectory: Directory[];
  isLoading: boolean;
  hasNewFolder: boolean;
  hasNewFile: boolean;
  hasFileRemoved: boolean;
  hasFileProtected: boolean;
  hasFileCutOrCopied: boolean;
  breadcrumb: MenuItem[];
  detailFolder: Directory[];
  currentPath: string;
  folderRequirePassword: boolean;
  searchResults: string[];
  hasFileRenamed: boolean;
  fileContent: string;
  imagesPool: ImageViewer[];
};

const initialState: StorageState = {
  status: null,
  allDirectorys: [],
  subMenuDirectory: [],
  isLoading: false,
  hasNewFolder: false,
  hasNewFile: false,
  hasFileRemoved: false,
  hasFileProtected: false,
  hasFileCutOrCopied: false,
  breadcrumb: [],
  detailFolder: [],
  currentPath: '',
  folderRequirePassword: false,
  searchResults: [],
  hasFileRenamed: false,
  fileContent: '',
  imagesPool: [],
};

export const StorageStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withComputed(({ status, allDirectorys }) => ({
    percentage: computed(() => {
      return ((status()?.used || 0) / (status()?.maximunSize || 0)) * 100;
    }),
    colorStatus: computed(() => {
      const percentage =
        ((status()?.used || 0) / (status()?.maximunSize || 0)) * 100;

      for (let i = 0; i < COLOR_STATUS_STORAGE.length; i++) {
        if (percentage <= COLOR_STATUS_STORAGE[i].threshold) {
          return COLOR_STATUS_STORAGE[i].color;
        }
      }
      return DEFAULT_STATUS_STORAGE_COLOR;
    }),
    directories: computed<TreeNode[]>(() => {
      if (allDirectorys() === null || allDirectorys().length === 0) {
        return [];
      }
      return allDirectorys().map((dir) => ({
        key: dir.name,
        label: dir.name,
        data: dir,
        leaf: false,
        loading: false,
        ...(dir.type === 'folder' && { children: [] }),
      }));
    }),
  })),
  withMethods(
    (
      store,
      storageService = inject(StorageService),
      messageService = inject(MessageService)
    ) => ({
      getSystemStorageStatus: rxMethod<void>(
        pipe(
          switchMap(() => {
            patchState(store, { isLoading: true });
            return storageService.getSystemStorageStatus().pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { status: res.response });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      getDirectory: rxMethod<string>(
        pipe(
          switchMap((location) => {
            patchState(store, { isLoading: true });
            return storageService.getDirectory(location).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { allDirectorys: res.response });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      getDetailsDirectory: rxMethod<{
        path: string;
        type: 'subMenu' | 'detailFolder';
        credential?: string;
      }>(
        pipe(
          switchMap((payload) => {
            if (payload.type === 'detailFolder') {
              patchState(store, {
                isLoading: true,
                currentPath: payload.path,
              });
            }
            return storageService
              .getDirectory(payload.path, payload.credential)
              .pipe(
                tapResponse({
                  next: (res) => {
                    patchState(store, { folderRequirePassword: false });
                    if (payload.type === 'subMenu') {
                      patchState(store, { subMenuDirectory: res.response });
                    } else {
                      const imagesPool = res.response
                        .filter((dir) =>
                          dir.extension
                            .toLowerCase()
                            .match(/(jpg|jpeg|png|gif)$/)
                        )
                        .map((dir) => {
                          return {
                            itemImageSrc: '',
                            thumbnailImageSrc: '',
                            alt: dir.name,
                            title: dir.name,
                          };
                        });
                      patchState(store, {
                        detailFolder: res.response,
                        imagesPool,
                      });
                    }
                  },
                  error: (err: HttpErrorResponse) => {
                    if (err.error.errorCode === 1018) {
                      patchState(store, { folderRequirePassword: true });
                    }
                  },
                  finalize: () => patchState(store, { isLoading: false }),
                })
              );
          })
        )
      ),
      uploadFile: rxMethod<UploadFileRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, hasNewFile: false });
            return storageService.uploadFile(payload).pipe(
              tapResponse({
                next: () => {
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'File uploaded successfully',
                  });
                  patchState(store, { hasNewFile: true });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      downloadFile: rxMethod<DownloadFileRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return storageService.downloadFile(payload).pipe(
              tapResponse({
                next: (res) => {
                  const blob = new Blob([res.response], {
                    type: 'application/octet-stream',
                  });
                  const url = window.URL.createObjectURL(blob);
                  const a = document.createElement('a');
                  a.href = url;
                  a.download = payload.request.fileNameToDownload;
                  a.click();
                  window.URL.revokeObjectURL(url);
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      createFolder: rxMethod<string>(
        pipe(
          switchMap((folderName) => {
            patchState(store, { isLoading: true });
            return storageService.createFolder(folderName).pipe(
              tapResponse({
                next: () => {
                  patchState(store, { hasNewFolder: true });
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Folder created successfully',
                  });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      removeFile: rxMethod<RemoveFileRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, hasFileRemoved: false });
            return storageService.removeFile(payload).pipe(
              tapResponse({
                next: () => {
                  patchState(store, { hasFileRemoved: true });
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'File removed successfully',
                  });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      folderProtection: rxMethod<FolderProtectionRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, hasFileProtected: false });
            return storageService.setFolderProtection(payload).pipe(
              tapResponse({
                next: () => {
                  patchState(store, { hasFileProtected: true });
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Folder protected successfully',
                  });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      checkFolderProtection: rxMethod<CheckFolderProtectionRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return storageService.checkFolderProtection(payload).pipe(
              tapResponse({
                next: () => {
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Folder protected successfully',
                  });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      searchFile: rxMethod<SearchRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return storageService.searchFile(payload).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { searchResults: res.response });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      renameFileOrFolder: rxMethod<RenameRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, hasFileRenamed: false });
            return storageService.rename(payload).pipe(
              tapResponse({
                next: () => {
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'File/Folder renamed successfully',
                  });
                  patchState(store, { hasFileRenamed: true });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () =>
                  patchState(store, {
                    isLoading: false,
                  }),
              })
            );
          })
        )
      ),
      readFileContent: rxMethod<ReadFileRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, fileContent: '' });
            return storageService.readFileContent(payload).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { fileContent: res });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
      cutOrCopy: rxMethod<CutOrCopyRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true, hasFileCutOrCopied: false });
            return storageService.cutOrCopy(payload).pipe(
              tapResponse({
                next: () => {
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'File/Folder copied successfully',
                  });
                  patchState(store, { hasFileCutOrCopied: true });
                },
                error: (err) => {
                  console.log(err);
                },
                finalize: () => patchState(store, { isLoading: false }),
              })
            );
          })
        )
      ),
    })
  )
);

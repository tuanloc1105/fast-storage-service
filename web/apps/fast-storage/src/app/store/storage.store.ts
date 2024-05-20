import { computed, inject } from '@angular/core';
import { StorageService } from '@app/data-access';
import {
  COLOR_STATUS_STORAGE,
  DEFAULT_STATUS_STORAGE_COLOR,
} from '@app/shared/constant';
import { Directory, DirectoryRequest, StorageStatus } from '@app/shared/model';
import { tapResponse } from '@ngrx/operators';
import {
  patchState,
  signalStore,
  withComputed,
  withMethods,
  withState,
} from '@ngrx/signals';
import { rxMethod } from '@ngrx/signals/rxjs-interop';
import { MessageService, TreeNode } from 'primeng/api';
import { pipe, switchMap } from 'rxjs';

type StorageState = {
  status: StorageStatus | null;
  directory: Directory[];
  isLoading: boolean;
};

const initialState: StorageState = {
  status: null,
  directory: [],
  isLoading: false,
};

export const StorageStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withComputed(({ status, directory }) => ({
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
      if (directory() === null || directory().length === 0) {
        return [];
      }
      return directory().map((dir, index) => ({
        key: index.toString(),
        label: dir.name,
        data: dir,
        icon: dir.type === 'folder' ? 'pi pi-fw pi-folder' : 'pi pi-fw pi-file',
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
      getDirectory: rxMethod<DirectoryRequest>(
        pipe(
          switchMap((payload) => {
            patchState(store, { isLoading: true });
            return storageService.getDirectory(payload).pipe(
              tapResponse({
                next: (res) => {
                  patchState(store, { directory: res.response });
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
      uploadFile: rxMethod<File>(
        pipe(
          switchMap((file) => {
            patchState(store, { isLoading: true });
            return storageService.uploadFile(file).pipe(
              tapResponse({
                next: (res) => {
                  messageService.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'File uploaded successfully',
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
      downloadFile: rxMethod<string>(
        pipe(
          switchMap((fileName) => {
            patchState(store, { isLoading: true });
            return storageService.downloadFile(fileName).pipe(
              tapResponse({
                next: (res) => {
                  const url = window.URL.createObjectURL(res.response);
                  const a = document.createElement('a');
                  a.href = url;
                  a.download = fileName;
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
    })
  )
);

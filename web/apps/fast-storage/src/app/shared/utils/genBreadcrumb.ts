import { patchState } from '@ngrx/signals';

export function genBreadcrumb(store: any, path: string[]): void {
  patchState(store, {
    currentPath: path.join('/'),
    breadcrumb: path.map((p) => ({
      label: p,
      command: () => {
        const nextPath = path.slice(0, path.indexOf(p) + 1).join('/');
        store.getDetailsDirectory({
          path: nextPath,
          type: 'detailFolder',
        });
      },
      styleClass: 'cursor-pointer',
    })),
  });
}

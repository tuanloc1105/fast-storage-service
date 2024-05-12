import { patchState, signalStore, withMethods, withState } from '@ngrx/signals';
import { TreeNode } from 'primeng/api';

type AppState = {
  isDarkMode: boolean;
  selectedFolder: TreeNode<any> | null;
};

const initialState: AppState = {
  isDarkMode: false,
  selectedFolder: null,
};

export const AppStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store) => ({
    toggleDarkMode() {
      patchState(store, (state) => ({ isDarkMode: !state.isDarkMode }));
    },
  }))
);

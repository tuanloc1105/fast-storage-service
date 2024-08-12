import { patchState, signalStore, withMethods, withState } from '@ngrx/signals';

type BannerState = {
  isBannerVisible: boolean;
  title: string;
  content: string;
  action: string;
};

const initialState: BannerState = {
  isBannerVisible: false,
  title: '',
  content: '',
  action: '',
};

export const BannerStore = signalStore(
  { providedIn: 'root' },
  withState(initialState),
  withMethods((store) => ({
    showBanner(title: string, content: string, action: string) {
      patchState(store, () => ({
        isBannerVisible: true,
        title,
        content,
        action,
      }));
    },
    hideBanner() {
      patchState(store, () => ({
        isBannerVisible: false,
        title: '',
        content: '',
        action: '',
      }));
    },
  }))
);

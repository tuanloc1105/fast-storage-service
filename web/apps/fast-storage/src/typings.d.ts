export interface IElectronAPI {
  getAppVersion: () => Promise<any>;
  log: (message: string) => Promise<any>;
}

declare global {
  interface Window {
    electronAPI: IElectronAPI;
  }
}

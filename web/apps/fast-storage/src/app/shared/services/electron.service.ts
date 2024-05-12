import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ElectronService {
  get isElectron(): boolean {
    return !!(window && window.electronAPI);
  }

  log(message: string) {
    if (this.isElectron) {
      window.electronAPI.log(message);
    } else {
      console.log(message);
    }
  }
}

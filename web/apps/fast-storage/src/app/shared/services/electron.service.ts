import { Injectable } from '@angular/core';
import { environment } from 'environments/environment';

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
    } else if (!environment.production) {
      console.log(message);
    }
  }
}

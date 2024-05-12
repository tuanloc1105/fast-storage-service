import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class StorageService {
  #http = inject(HttpClient);

  constructor() {}

  public getSystemStorageStatus(): Observable<any> {
    return this.#http.get('/storage/system_storage_status');
  }
}

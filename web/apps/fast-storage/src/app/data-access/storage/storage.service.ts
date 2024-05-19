import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Directory, StorageStatus } from '@app/shared/model';
import { CommonResponse } from '@app/shared/model/common.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class StorageService {
  private readonly http = inject(HttpClient);

  public getSystemStorageStatus(): Observable<CommonResponse<StorageStatus>> {
    return this.http.get<CommonResponse<StorageStatus>>(
      '/storage/user_storage_status'
    );
  }

  public getDirectory(): Observable<CommonResponse<Directory[]>> {
    return this.http.post<CommonResponse<Directory[]>>(
      '/storage/get_all_element_in_specific_directory',
      {}
    );
  }

  public uploadFile(file: File): Observable<CommonResponse<any>> {
    const formData = new FormData();
    formData.append('file', file);
    return this.http.post<CommonResponse<any>>(
      '/storage/upload_file',
      formData
    );
  }

  public downloadFile(fileName: string): Observable<CommonResponse<any>> {
    return this.http.get<CommonResponse<any>>(
      `/storage/download_file/${fileName}`
    );
  }
}

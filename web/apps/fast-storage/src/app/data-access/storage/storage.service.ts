import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import {
  CreateFolderRequest,
  Directory,
  DirectoryRequest,
  StorageStatus,
  UploadFileRequest,
} from '@app/shared/model';
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

  public getDirectory(
    location: string
  ): Observable<CommonResponse<Directory[]>> {
    const payload: DirectoryRequest = {
      request: {
        currentLocation: location,
      },
    };
    return this.http.post<CommonResponse<Directory[]>>(
      '/storage/get_all_element_in_specific_directory',
      payload
    );
  }

  public uploadFile(
    payload: UploadFileRequest
  ): Observable<CommonResponse<any>> {
    const formData = new FormData();
    formData.append('file', payload.file);
    formData.append('folderLocation', payload.folderLocation);
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

  public createFolder(folderName: string): Observable<CommonResponse<any>> {
    const payload: CreateFolderRequest = {
      request: {
        folderToCreate: folderName,
      },
    };
    return this.http.post<CommonResponse<any>>(
      '/storage/create_folder',
      payload
    );
  }
}

/* eslint-disable @typescript-eslint/no-explicit-any */
import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import {
  CheckFolderProtectionRequest,
  CreateFolderRequest,
  CutOrCopyRequest,
  Directory,
  DirectoryRequest,
  DownloadFileRequest,
  FolderProtectionRequest,
  ReadFileRequest,
  RemoveFileRequest,
  RenameRequest,
  SearchRequest,
  ShowImageRequest,
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
    location: string,
    credential?: string
  ): Observable<CommonResponse<Directory[]>> {
    const payload: DirectoryRequest = {
      request: {
        currentLocation: location,
        credential,
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
    payload.files.forEach((file) => {
      formData.append('file', file);
    });
    formData.append('folderLocation', payload.folderLocation);
    return this.http.post<CommonResponse<any>>(
      '/storage/upload_file',
      formData
    );
  }

  public downloadFile(
    payload: DownloadFileRequest
  ): Observable<CommonResponse<any>> {
    const { fileNameToDownload, locationToDownload } = payload.request;

    return this.http.get<CommonResponse<any>>(
      `/storage/download_file?fileNameToDownload=${fileNameToDownload}&locationToDownload=${locationToDownload}`
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

  public removeFile(
    payload: RemoveFileRequest
  ): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>('/storage/remove_file', payload);
  }

  public setFolderProtection(
    payload: FolderProtectionRequest
  ): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>(
      '/storage/set_password_for_folder',
      payload
    );
  }

  public checkFolderProtection(
    payload: CheckFolderProtectionRequest
  ): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>(
      '/storage/check_secure_folder_status',
      payload
    );
  }

  public searchFile(payload: SearchRequest): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>('/storage/search_file', payload);
  }

  public rename(payload: RenameRequest): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>(
      '/storage/rename_file_or_folder',
      payload
    );
  }

  public readFileContent(payload: ReadFileRequest): Observable<string> {
    const { fileNameToRead, locationToRead } = payload;
    return this.http.get('/storage/read_text_file_content', {
      responseType: 'text',
      params: {
        fileNameToRead,
        locationToRead,
      },
    });
  }

  public showImage(
    payload: ShowImageRequest
  ): Observable<CommonResponse<{ data: string }>> {
    return this.http.post<CommonResponse<{ data: string }>>(
      '/storage/read_image_file',
      payload
    );
  }

  public cutOrCopy(payload: CutOrCopyRequest): Observable<CommonResponse<any>> {
    return this.http.post<CommonResponse<any>>('/storage/cut_or_copy', payload);
  }
}

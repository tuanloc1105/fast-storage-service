export interface StorageStatus {
  maximunSize: number;
  used: number;
}

export interface Directory {
  name: string;
  size: string;
  type: 'file' | 'folder';
  extension: string;
  lastModifiedDate: string;
}

export interface DirectoryRequest {
  request: {
    currentLocation: string;
  };
}

export interface CreateFolderRequest {
  request: {
    folderToCreate: string;
  };
}

export interface UploadFileRequest {
  files: File[];
  folderLocation: string;
}

export interface DownloadFileRequest {
  request: {
    fileNameToDownload: string;
    locationToDownload: string;
  };
}

export interface RemoveFileRequest {
  request: {
    fileNameToRemove: string;
    locationToRemove: string;
    otpCredential?: string;
  };
}

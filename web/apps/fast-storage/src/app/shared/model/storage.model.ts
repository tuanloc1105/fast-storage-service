export interface StorageStatus {
  maximunSize: number;
  used: number;
}

export interface Directory {
  name: string;
  size: string;
  type: string;
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
  file: File;
  folderLocation: string;
}

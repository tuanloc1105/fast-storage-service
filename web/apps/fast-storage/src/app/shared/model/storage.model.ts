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
    credential?: string;
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

export interface FolderProtectionRequest {
  request: {
    folder: string;
    credentialType: 'PASSWORD' | 'OTP';
    credential: string;
  };
}

export interface CheckFolderProtectionRequest {
  request: {
    folder: string;
  };
}

export interface SearchRequest {
  request: {
    searchingContent: string;
  };
}

export interface RenameRequest {
  request: {
    oldFolderLocationName: string;
    newFolderLocationName: string;
  };
}

export interface ReadFileRequest {
  fileNameToRead: string;
  locationToRead: string;
}

export interface ShowImageRequest {
  request: {
    folderLocation: string;
    imageFileName: string;
  };
}

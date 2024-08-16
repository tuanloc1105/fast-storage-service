export interface CommonResponse<T> {
  errorCode: number;
  errorMessage: string;
  trace: string;
  response: T;
}

export interface ImageViewer {
  itemImageSrc: string;
  thumbnailImageSrc: string;
  alt: string;
  title: string;
}

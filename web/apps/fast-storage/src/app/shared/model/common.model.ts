export interface CommonResponse<T> {
  errorCode: number;
  errorMessage: string;
  trace: string;
  response: T;
}

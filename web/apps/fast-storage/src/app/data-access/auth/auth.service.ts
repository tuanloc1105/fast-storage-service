import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import {
  GetNewTokenRequest,
  LoginRequest,
  LoginResponse,
  LogoutRequest,
  RegisterRequest,
  UserInfoResponse,
} from '@app/shared/model';
import { CommonResponse } from '@app/shared/model/common.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly http = inject(HttpClient);

  public getUserInfo(): Observable<CommonResponse<UserInfoResponse>> {
    return this.http.get<CommonResponse<UserInfoResponse>>(
      '/auth/get_user_info'
    );
  }

  public login(
    payload: LoginRequest
  ): Observable<CommonResponse<LoginResponse>> {
    return this.http.post<CommonResponse<LoginResponse>>(
      '/auth/login',
      payload
    );
  }

  public logout(payload: LogoutRequest): Observable<any> {
    return this.http.post('/auth/logout', payload);
  }

  public register(payload: RegisterRequest): Observable<any> {
    return this.http.post('/auth/register', payload);
  }

  public getNewToken(payload: GetNewTokenRequest): Observable<any> {
    return this.http.post('/auth/get_new_token', payload);
  }
}

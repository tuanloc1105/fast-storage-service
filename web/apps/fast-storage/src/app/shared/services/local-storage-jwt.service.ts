import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class LocalStorageJwtService {
  getAccessToken(): Observable<string | null> {
    const data = localStorage.getItem('access_token');
    if (data) {
      return of(data);
    }
    return of(null);
  }

  getRefreshToken(): Observable<string | null> {
    const data = localStorage.getItem('refresh_token');
    if (data) {
      return of(data);
    }
    return of(null);
  }

  setItem(data: {
    access_token: string;
    refresh_token: string;
  }): Observable<boolean> {
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);
    return of(true);
  }

  removeItem(): Observable<boolean> {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    return of(true);
  }
}

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

  getIsFirstTime(): { is: boolean; date: string } | null {
    return JSON.parse(localStorage.getItem('isFirstTime') || '{}');
  }

  setItem(data: {
    access_token: string;
    refresh_token: string;
  }): Observable<boolean> {
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);

    const isFirstTime = localStorage.getItem('isFirstTime');
    if (!isFirstTime) {
      const data = { is: true, date: new Date().toISOString() };
      localStorage.setItem('isFirstTime', JSON.stringify(data));
    } else if (new Date().getDay() > new Date(isFirstTime).getDay()) {
      const data = { is: true, date: new Date().toISOString() };
      localStorage.setItem('isFirstTime', JSON.stringify(data));
    } else {
      const data = { is: false, date: new Date().toISOString() };
      localStorage.setItem('isFirstTime', JSON.stringify(data));
    }
    return of(true);
  }

  removeItem(): Observable<boolean> {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    return of(true);
  }
}

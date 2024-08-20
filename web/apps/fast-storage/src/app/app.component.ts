import { DOCUMENT } from '@angular/common';
import { Component, Inject, OnInit, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { patchState } from '@ngrx/signals';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { PrimeNGConfig } from 'primeng/api';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { ToastModule } from 'primeng/toast';
import { lastValueFrom } from 'rxjs';
import { LocalStorageJwtService } from './shared/services';
import { AppStore, StorageStore } from './store';

@Component({
  standalone: true,
  imports: [RouterOutlet, TranslateModule, ToastModule, ConfirmDialogModule],
  selector: 'app-root',
  template: '<p-toast /><p-confirmDialog /><router-outlet />',
  styles: [],
})
export class AppComponent implements OnInit {
  private readonly primengConfig = inject(PrimeNGConfig);
  private readonly localStorageJwtService = inject(LocalStorageJwtService);
  private readonly router = inject(Router);
  private readonly appStore = inject(AppStore);
  private readonly storageStore = inject(StorageStore);

  constructor(
    translate: TranslateService,
    @Inject(DOCUMENT) private document: Document
  ) {
    // this language will be used as a fallback when a translation isn't found in the current language
    translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    translate.use('en');
  }

  ngOnInit() {
    this.primengConfig.ripple = true;
    this.initApp();
  }

  private async initApp() {
    const accessToken = await lastValueFrom(
      this.localStorageJwtService.getAccessToken()
    );

    const currentTheme = localStorage.getItem('theme');
    const themeLink = this.document.getElementById(
      'app-theme'
    ) as HTMLLinkElement;
    if (currentTheme) {
      if (currentTheme === 'dark') {
        themeLink.href = 'aura-dark-cyan.css';
        patchState(this.appStore, { isDarkMode: true });
      } else {
        patchState(this.appStore, { isDarkMode: false });
        themeLink.href = 'aura-light-cyan.css';
      }
    } else {
      localStorage.setItem('theme', 'light');
      patchState(this.appStore, { isDarkMode: false });
      themeLink.href = 'aura-light-cyan.css';
    }

    if (accessToken) {
      const queryString = window.location.search;
      if (queryString.includes('path')) {
        const path = new URLSearchParams(queryString).get('path') || '';
        patchState(this.storageStore, { currentPath: path });
      }

      this.router.navigate(['app/initializing']);
    }
  }
}

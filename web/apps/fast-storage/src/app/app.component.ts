import { Component, Inject, OnInit, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { PrimeNGConfig } from 'primeng/api';
import { lastValueFrom } from 'rxjs';
import { LocalStorageJwtService } from './shared/services';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { DOCUMENT } from '@angular/common';
import { AppStore } from './store';
import { patchState } from '@ngrx/signals';

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
      this.router.navigate([
        'app/initializing',
        { returnUrl: location.pathname },
      ]);
    }
  }
}

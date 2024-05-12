import { Component, OnInit, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { PrimeNGConfig } from 'primeng/api';
import { lastValueFrom } from 'rxjs';
import { LocalStorageJwtService } from './shared/services';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';

@Component({
  standalone: true,
  imports: [RouterOutlet, TranslateModule, ToastModule, ConfirmDialogModule],
  selector: 'app-root',
  template: '<p-toast /><p-confirmDialog /><router-outlet />',
  styles: [],
})
export class AppComponent implements OnInit {
  #primengConfig = inject(PrimeNGConfig);
  #localStorageJwtService = inject(LocalStorageJwtService);
  #router = inject(Router);

  constructor(translate: TranslateService) {
    // this language will be used as a fallback when a translation isn't found in the current language
    translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    translate.use('en');
  }

  ngOnInit() {
    this.#primengConfig.ripple = true;
    this.initApp();
  }

  private async initApp() {
    const accessToken = await lastValueFrom(
      this.#localStorageJwtService.getAccessToken()
    );

    if (accessToken) {
      this.#router.navigate(['app/initializing']);
    } else {
      this.#router.navigate(['auth/login']);
    }
  }
}

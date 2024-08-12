import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { SlideInDirective } from '@app/shared/directives';
import { BannerStore } from '@app/store';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-banner',
  standalone: true,
  imports: [CommonModule, ButtonModule, SlideInDirective],
  templateUrl: './banner.component.html',
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class BannerComponent {
  public bannerStore = inject(BannerStore);

  public hide() {
    sessionStorage.setItem('offerBanner', 'false');
    this.bannerStore.hideBanner();
  }
}

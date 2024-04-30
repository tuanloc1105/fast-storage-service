import { JsonPipe } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  effect,
  inject,
} from '@angular/core';
import { AppStore } from '@app/state/app.state';
import { getState } from '@ngrx/signals';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';

@Component({
  selector: 'app-folder-detail',
  standalone: true,
  imports: [ButtonModule, JsonPipe, TableModule],
  templateUrl: './folder-detail.component.html',
  styleUrl: './folder-detail.component.scss',
  providers: [],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderDetailComponent {
  readonly appStore = inject(AppStore);

  constructor() {
    effect(() => {
      const state = getState(this.appStore);
      console.log('FolderDetailComponent', state);
    });
  }
}

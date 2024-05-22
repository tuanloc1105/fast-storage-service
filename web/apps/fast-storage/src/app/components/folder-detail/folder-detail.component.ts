import { JsonPipe } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  OnInit,
  effect,
  inject,
} from '@angular/core';
import { Directory } from '@app/shared/model';
import { StorageStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { MenuItem } from 'primeng/api';
import { BreadcrumbItemClickEvent, BreadcrumbModule } from 'primeng/breadcrumb';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';

@Component({
  selector: 'app-folder-detail',
  standalone: true,
  imports: [ButtonModule, JsonPipe, BreadcrumbModule, TableModule],
  templateUrl: './folder-detail.component.html',
  styleUrl: './folder-detail.component.scss',
  providers: [],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderDetailComponent implements OnInit {
  public storageStore = inject(StorageStore);

  public home: MenuItem | undefined;
  public selectedDirectory: Directory | null = null;

  constructor() {
    effect(() => {
      console.log(this.storageStore.breadcrumb());
    });
  }

  ngOnInit(): void {
    this.home = { icon: 'pi pi-home' };
  }

  public handleBreadcrumb(event: BreadcrumbItemClickEvent): void {
    const item = event.item.label;
    const index = this.storageStore
      .breadcrumb()
      .findIndex((b) => b.label === item);
    patchState(this.storageStore, {
      breadcrumb: this.storageStore.breadcrumb().slice(0, index + 1),
    });
  }
}

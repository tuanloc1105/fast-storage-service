import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  computed,
  DestroyRef,
  inject,
  OnInit,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { genBreadcrumb } from '@app/shared/utils';
import { StorageStore } from '@app/store';
import { DialogModule } from 'primeng/dialog';
import { DividerModule } from 'primeng/divider';
import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { InputTextModule } from 'primeng/inputtext';
import { ProgressBarModule } from 'primeng/progressbar';
import { Subject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    CommonModule,
    DialogModule,
    InputIconModule,
    InputTextModule,
    FormsModule,
    DividerModule,
    IconFieldModule,
    ProgressBarModule,
  ],
  template: `<div class="flex flex-col gap-2">
    <p-iconField iconPosition="left">
      <p-inputIcon styleClass="pi pi-search" />
      <input
        pInputText
        icon="pi pi-search"
        type="text"
        [(ngModel)]="searchValue"
        (ngModelChange)="searchSubject.next($event)"
        placeholder="Search..."
        class="border-none rounded-lg outline-none w-full"
      />
    </p-iconField>
    @if (storageStore.isLoading()) {
    <p-progressBar mode="indeterminate" [style]="{ height: '6px' }" />
    }
    <p-divider />
    @if (storageStore.searchResults().length > 0) {
    <div class="flex flex-col gap-2">
      <div class="flex flex-col gap-3">
        @for (item of storageStore.searchResults(); track $index) {
        <div
          class="group flex items-center justify-between p-4 bg-[#09090b] rounded-lg cursor-pointer hover:bg-blue-400"
          (click)="viewFile(item)"
        >
          <span>{{ item }}</span>
          <i
            class="pi pi-arrow-right group-hover:translate-x-1 transition ease-in-out delay-150"
          ></i>
        </div>
        }
      </div>
    </div>
    } @else {
    <span class="text-center p-12">No file found using that search term.</span>
    }
  </div>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent implements OnInit {
  public storageStore = inject(StorageStore);

  private ref = inject(DynamicDialogRef);
  private destroyRef = inject(DestroyRef);

  public searchSubject = new Subject<string>();
  private readonly debounceTimeMs = 500;

  public searchValue = '';

  ngOnInit(): void {
    this.searchSubject
      .pipe(
        debounceTime(this.debounceTimeMs),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((searchValue) => {
        this.handleSearch(searchValue);
      });
  }

  public handleSearch(value: string) {
    this.storageStore.searchFile({ request: { searchingContent: value } });
  }

  public viewFile(path: string) {
    const folderPath = path.split('/').slice(0, -1).join('/');
    this.storageStore.getDetailsDirectory({
      path: folderPath,
      type: 'detailFolder',
    });
    genBreadcrumb(this.storageStore, path.split('/').slice(0, -1));
    this.ref.close();
  }
}

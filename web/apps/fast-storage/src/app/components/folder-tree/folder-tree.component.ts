import { DOCUMENT } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  computed,
  effect,
  inject,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NewFolderComponent } from '@app/shared/components';
import { Directory } from '@app/shared/model';
import { getFullPath } from '@app/shared/utils';
import { AppStore, StorageStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { TreeNode } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { DialogService, DynamicDialogModule } from 'primeng/dynamicdialog';
import { InputTextModule } from 'primeng/inputtext';
import { MeterGroupModule } from 'primeng/metergroup';
import { PanelMenuModule } from 'primeng/panelmenu';
import {
  TreeModule,
  TreeNodeCollapseEvent,
  TreeNodeExpandEvent,
} from 'primeng/tree';

@Component({
  selector: 'app-folder-tree',
  standalone: true,
  imports: [
    ButtonModule,
    MeterGroupModule,
    PanelMenuModule,
    TreeModule,
    DialogModule,
    InputTextModule,
    FormsModule,
    DynamicDialogModule,
  ],
  templateUrl: './folder-tree.component.html',
  styleUrl: './folder-tree.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderTreeComponent implements OnInit {
  public appStore = inject(AppStore);
  public storageStore = inject(StorageStore);

  private readonly dialogService = inject(DialogService);

  public selectedDocumentFolder: TreeNode | null = null;

  private nodeExpandEvent!: TreeNodeExpandEvent;
  private cdr = inject(ChangeDetectorRef);

  public meter = computed(() => [
    {
      label: 'Space used',
      value: this.storageStore.percentage(),
      color: this.storageStore.colorStatus(),
    },
  ]);

  constructor(@Inject(DOCUMENT) private document: Document) {
    effect(
      () => {
        if (
          this.storageStore.subMenuDirectory() &&
          this.storageStore.subMenuDirectory().length > 0
        ) {
          this.handleRefreshFolder(this.storageStore.subMenuDirectory());
        } else {
          if (this.nodeExpandEvent) {
            this.nodeExpandEvent.node.loading = false;
            this.cdr.detectChanges();
          }
        }
      },
      { allowSignalWrites: true }
    );
  }

  ngOnInit(): void {
    this.storageStore.getDirectory('');
  }

  public switchTheme() {
    const currentTheme = localStorage.getItem('theme');
    const themeLink = this.document.getElementById(
      'app-theme'
    ) as HTMLLinkElement;
    if (currentTheme === 'dark') {
      themeLink.href = 'aura-light-cyan.css';
      patchState(this.appStore, { isDarkMode: false });
      localStorage.setItem('theme', 'light');
    } else {
      themeLink.href = 'aura-dark-cyan.css';
      patchState(this.appStore, { isDarkMode: true });
      localStorage.setItem('theme', 'dark');
    }
  }

  public onNodeExpand(event: TreeNodeExpandEvent) {
    this.nodeExpandEvent = event;
    this.nodeExpandEvent.node.loading = true;
    const path = getFullPath(event.node);
    this.storageStore.getDetailsDirectory({ path, type: 'subMenu' });
  }

  public onNodeCollapse(event: TreeNodeCollapseEvent) {
    event.node.children = [];
  }

  public addNewFolder() {
    if (this.selectedDocumentFolder) {
      const path = getFullPath(this.selectedDocumentFolder);
      patchState(this.storageStore, { currentPath: path });
    }

    this.dialogService.open(NewFolderComponent, {
      header: 'Create new folder',
    });
  }

  public onNodeSelect(event: TreeNode<any> | TreeNode<any>[] | null) {
    if (event && !Array.isArray(event)) {
      const path = getFullPath(event).split('/');
      this.storageStore.getDetailsDirectory({
        path: path.join('/'),
        type: 'detailFolder',
      });
      patchState(this.storageStore, {
        currentPath: path.join('/'),
        breadcrumb: path.map((p) => ({
          label: p,
          command: () => {
            const nextPath = path.slice(0, path.indexOf(p) + 1).join('/');
            this.storageStore.getDetailsDirectory({
              path: nextPath,
              type: 'detailFolder',
            });
          },
          styleClass: 'cursor-pointer',
        })),
      });
    }
  }

  private handleRefreshFolder(data: Directory[]) {
    const _node = { ...this.nodeExpandEvent?.node };
    const path = getFullPath(_node).split('/');

    data.forEach((dir) => {
      if (dir.type !== 'file')
        _node.children?.push({
          key: dir.name,
          label: dir.name,
          data: dir,
          leaf: false,
          loading: false,
          children: [],
        });
    });

    this.addChildrenToNode(this.storageStore.directories(), _node, path);
    this.nodeExpandEvent.node.loading = false;
    this.cdr.detectChanges();
  }

  private addChildrenToNode(
    directories: TreeNode[],
    node: TreeNode,
    path: string[]
  ) {
    directories.forEach((dir) => {
      if (dir.key === path[0]) {
        if (path.length === 1) {
          dir.children = node.children;
        } else {
          this.addChildrenToNode(dir.children || [], node, path.slice(1));
        }
      }
    });
  }
}

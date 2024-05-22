import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  computed,
  effect,
  inject,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Directory } from '@app/shared/model';
import { getFullPath, getNestedObject } from '@app/shared/utils';
import { AppStore, StorageStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { MenuItem, TreeNode } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { ContextMenuModule } from 'primeng/contextmenu';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { MeterGroupModule } from 'primeng/metergroup';
import { PanelMenuModule } from 'primeng/panelmenu';
import {
  TreeModule,
  TreeNodeContextMenuSelectEvent,
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
    ContextMenuModule,
  ],
  templateUrl: './folder-tree.component.html',
  styleUrl: './folder-tree.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderTreeComponent implements OnInit {
  public appStore = inject(AppStore);
  public storageStore = inject(StorageStore);

  public newFolderVisible = false;
  public newFolderName = '';

  public selectedDocumentFolder: TreeNode | null = null;
  public folderContextMenu: MenuItem[] = [];

  private nodeExpandEvent!: TreeNodeExpandEvent;
  private cdr = inject(ChangeDetectorRef);

  public meter = computed(() => [
    {
      label: 'Space used',
      value: this.storageStore.percentage(),
      color: this.storageStore.colorStatus(),
    },
  ]);

  constructor() {
    effect(
      () => {
        if (this.storageStore.hasNewFolder()) {
          this.closeNewFolderDialog();
        }

        if (
          this.storageStore.subMenuDirectory() &&
          this.storageStore.subMenuDirectory().length > 0
        ) {
          this.handleRefreshFolder(this.storageStore.subMenuDirectory());
        } else {
          if (this.nodeExpandEvent) this.nodeExpandEvent.node.loading = false;
        }
      },
      { allowSignalWrites: true }
    );
  }

  ngOnInit(): void {
    this.storageStore.getDirectory('');
  }

  public onTreeContextMenuSelect(event: TreeNodeContextMenuSelectEvent) {
    if (event.node.data.type === 'file') {
      this.folderContextMenu = [
        {
          label: 'Download',
          icon: 'pi pi-download',
        },
      ];
    } else {
      this.folderContextMenu = [
        {
          label: 'Create Folder',
          icon: 'pi pi-folder',
          command: () => (this.newFolderVisible = true),
        },
        {
          label: 'Create File',
          icon: 'pi pi-file',
          command: () => (this.newFolderVisible = true),
        },
      ];
    }
  }

  public onNodeExpand(event: TreeNodeExpandEvent) {
    this.nodeExpandEvent = event;
    this.nodeExpandEvent.node.loading = true;
    const path = getFullPath(event.node);
    this.storageStore.getDetailsDirectory({ path, type: 'subMenu' });
  }

  public addNewFolder() {
    if (this.selectedDocumentFolder) {
      const path = getFullPath(this.selectedDocumentFolder);
      this.storageStore.createFolder(path + '/' + this.newFolderName);
    } else {
      this.storageStore.createFolder(this.newFolderName);
    }
  }

  public closeNewFolderDialog() {
    this.newFolderVisible = false;
    this.newFolderName = '';
  }

  public onNodeSelect(event: TreeNode<any> | TreeNode<any>[] | null) {
    if (event && !Array.isArray(event)) {
      const path = getFullPath(event).split('/');
      this.storageStore.getDetailsDirectory({
        path: path.join('/'),
        type: 'detailFolder',
      });
      patchState(this.storageStore, {
        breadcrumb: path.map((p) => ({
          label: p,
          command: () => {
            const newPath = path.slice(0, path.indexOf(p) + 1).join('/');
            this.storageStore.getDetailsDirectory({
              path: newPath,
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

    _node.children = [];

    data.forEach((dir) => {
      _node.children?.push({
        key: dir.name,
        label: dir.name,
        data: dir,
        icon: dir.type === 'folder' ? 'pi pi-fw pi-folder' : 'pi pi-fw pi-file',
        leaf: false,
        loading: false,
        ...(dir.type === 'folder' && { children: [] }),
      });
    });

    this.addChildrenToNode(this.storageStore.directories(), _node, path);

    this.nodeExpandEvent.node.loading = false;
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

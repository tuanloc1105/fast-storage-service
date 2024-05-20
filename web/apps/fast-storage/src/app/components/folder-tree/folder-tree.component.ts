import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  computed,
  inject,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AppStore, StorageStore } from '@app/store';
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

  private cd = inject(ChangeDetectorRef);

  public meter = computed(() => [
    {
      label: 'Space used',
      value: this.storageStore.percentage(),
      color: this.storageStore.colorStatus(),
    },
  ]);

  ngOnInit(): void {
    this.storageStore.getDirectory({ request: { currentLocation: '' } });
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
          command: () => this.addNewFolder(),
        },
        {
          label: 'Create File',
          icon: 'pi pi-file',
          command: () => this.addNewFolder(),
        },
      ];
    }
  }

  public onNodeExpand(event: TreeNodeExpandEvent) {
    console.log(event);
  }

  public addNewFolder() {
    this.storageStore.getDirectory({
      request: { currentLocation: this.newFolderName },
    });
  }
}

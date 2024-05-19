import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  computed,
  inject,
} from '@angular/core';
import { AppStore, StorageStore } from '@app/store';
import { patchState } from '@ngrx/signals';
import { TreeNode } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { MeterGroupModule } from 'primeng/metergroup';
import { PanelMenuModule } from 'primeng/panelmenu';
import { TreeModule } from 'primeng/tree';

@Component({
  selector: 'app-folder-tree',
  standalone: true,
  imports: [ButtonModule, MeterGroupModule, PanelMenuModule, TreeModule],
  templateUrl: './folder-tree.component.html',
  styleUrl: './folder-tree.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderTreeComponent implements OnInit {
  public appStore = inject(AppStore);
  public storageStore = inject(StorageStore);

  private cd = inject(ChangeDetectorRef);

  public meter = computed(() => [
    {
      label: 'Space used',
      value: this.storageStore.percentage(),
      color: this.storageStore.colorStatus(),
    },
  ]);

  ngOnInit(): void {
    this.storageStore.getDirectory();
  }

  public setSelectedFolder(folder: TreeNode<any> | TreeNode<any>[] | null) {
    if (typeof folder === 'object' && folder !== null) {
      patchState(this.appStore, { selectedFolder: folder as TreeNode<any> });
    }
  }

  public onNodeExpand(event: any) {
    if (!event.node.children) {
      event.node.loading = true;

      setTimeout(() => {
        const _node = { ...event.node };
        _node.children = [];

        for (let i = 0; i < 3; i++) {
          _node.children.push({
            key: event.node.key + '-' + i,
            label: 'Lazy ' + event.node.label + '-' + i,
          });
        }

        const key = parseInt(_node.key, 10);
        this.storageStore.directories()[key] = { ..._node, loading: false };
        this.cd.markForCheck();
      }, 500);
    }
  }
}

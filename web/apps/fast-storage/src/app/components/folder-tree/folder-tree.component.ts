import {
  ChangeDetectionStrategy,
  Component,
  inject,
  signal,
} from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { MeterGroupModule } from 'primeng/metergroup';
import { PanelMenuModule } from 'primeng/panelmenu';
import { TreeModule } from 'primeng/tree';
import { FolderTreeStore } from './folder-tree.state';
import { AppStore } from '@app/store';
import { TreeNode } from 'primeng/api';
import { patchState } from '@ngrx/signals';

@Component({
  selector: 'app-folder-tree',
  standalone: true,
  imports: [ButtonModule, MeterGroupModule, PanelMenuModule, TreeModule],
  templateUrl: './folder-tree.component.html',
  styleUrl: './folder-tree.component.scss',
  providers: [FolderTreeStore],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FolderTreeComponent {
  readonly folderTreeStore = inject(FolderTreeStore);
  readonly appStore = inject(AppStore);

  public value = signal([{ label: 'Space used', value: 15, color: '#34d399' }]);

  public setSelectedFolder(folder: TreeNode<any> | TreeNode<any>[] | null) {
    if (typeof folder === 'object' && folder !== null) {
      patchState(this.appStore, { selectedFolder: folder as TreeNode<any> });
    }
  }
}

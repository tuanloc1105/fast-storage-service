import { Component } from '@angular/core';
import {
  FolderDetailComponent,
  FolderTreeComponent,
  SidebarComponent,
} from '@app/components';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [SidebarComponent, FolderTreeComponent, FolderDetailComponent],
  template: `<div class="flex flex-row h-screen">
    <div class="basis-[5%] px-2 py-7 border-r">
      <app-sidebar></app-sidebar>
    </div>
    <div class="basis-1/3 px-4 py-7 border-r">
      <app-folder-tree></app-folder-tree>
    </div>
    <div class="basis-full px-5 py-7">
      <app-folder-detail></app-folder-detail>
    </div>
  </div>`,
  styles: [],
})
export class LayoutComponent {}

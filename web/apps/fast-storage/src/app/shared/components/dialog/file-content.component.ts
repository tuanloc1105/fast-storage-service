import {
  ChangeDetectionStrategy,
  Component,
  inject,
  OnInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { DynamicDialogConfig, DynamicDialogRef } from 'primeng/dynamicdialog';
import { MarkdownModule } from 'ngx-markdown';

@Component({
  selector: 'app-file-content',
  standalone: true,
  imports: [CommonModule, MarkdownModule],
  template: `<div>
    <markdown lineNumbers clipboard [data]="config.data.content"></markdown>
  </div>`,
  styles: ``,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FileContentComponent implements OnInit {
  private ref = inject(DynamicDialogRef);
  public config = inject(DynamicDialogConfig);

  ngOnInit(): void {
    console.log(this.config.data.content);
  }
}
